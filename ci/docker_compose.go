package main

import (
	"fmt"
	"sort"

	"dagger.io/dagger"
	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	*dagger.Container
	Name string
}

type PublishedPort struct {
	Address  string
	Target   int
	Protocol dagger.NetworkProtocol
}

func BindDockerCompose(c *dagger.Client, bindContainer *dagger.Container, dockerComposeFileName string) *dagger.Container {
	var services []*Service
	opts, err := cli.NewProjectOptions([]string{dockerComposeFileName},
		cli.WithWorkingDirectory("."),
		cli.WithDefaultConfigPath,
		cli.WithOsEnv,
		cli.WithConfigFileEnv,
	)
	if err != nil {
		log.Fatalf("Problem loading docker compose options %v", err)
	}

	project, err := cli.ProjectFromOptions(opts)
	if err != nil {
		log.Fatalf("Problem loading docker compose %v", err)
	}

	for _, svc := range project.Services {
		// Append the service to the list of services
		daggerSvc, err := serviceContainer(c, project, svc)
		if err != nil {
			log.Fatal(err)
		}
		services = append(services, daggerSvc)
	}
	for _, svc := range services {
		log.Infof("Bound with :alias %s", svc.Name)
		bindContainer = bindContainer.WithServiceBinding(svc.Name, svc.Container.AsService())
	}
	return bindContainer
}

func serviceContainer(c *dagger.Client, project *types.Project, svc types.ServiceConfig) (*Service, error) {
	ctr := c.Pipeline(svc.Name).Container(dagger.ContainerOpts{Platform: "linux/amd64"})
	if svc.Image != "" {
		ctr = ctr.From(svc.Image)
	} else if svc.Build != nil {
		args := []dagger.BuildArg{}
		for name, val := range svc.Build.Args {
			if val != nil {
				args = append(args, dagger.BuildArg{
					Name:  name,
					Value: *val,
				})
			}
		}

		ctr = ctr.Build(c.Host().Directory(svc.Build.Context), dagger.ContainerBuildOpts{
			Dockerfile: svc.Build.Dockerfile,
			BuildArgs:  args,
			Target:     svc.Build.Target,
		})
	}

	// sort env to ensure same container
	type env struct{ name, value string }
	envs := []env{}
	for name, val := range svc.Environment {
		if val != nil {
			envs = append(envs, env{name, *val})
		}
	}
	sort.Slice(envs, func(i, j int) bool {
		return envs[i].name < envs[j].name
	})
	for _, env := range envs {
		ctr = ctr.WithEnvVariable(env.name, env.value)
	}

	for _, port := range svc.Ports {
		log.Infof("Port: %v", port)
		switch port.Mode {
		case "ingress":
			ctr = ctr.WithExposedPort(int(port.Target))
		default:
			return nil, fmt.Errorf("port mode %s not supported", port.Mode)
		}
	}

	for _, vol := range svc.Volumes {
		switch vol.Type {
		case types.VolumeTypeBind:
			ctr = ctr.WithMountedDirectory(vol.Target, c.Host().Directory(vol.Source))
		case types.VolumeTypeVolume:
			ctr = ctr.WithMountedCache(vol.Target, c.CacheVolume(vol.Source))
		case types.VolumeTypeTmpfs:
			ctr = ctr.WithMountedTemp(vol.Target)
		default:
			return nil, fmt.Errorf("volume type %s not supported", vol.Type)
		}
	}

	for depName := range svc.DependsOn {
		cfg, err := project.GetService(depName)
		if err != nil {
			return nil, err
		}

		svc, err := serviceContainer(c, project, cfg)
		if err != nil {
			return nil, err
		}

		ctr = ctr.WithServiceBinding(depName, svc.Container.AsService())
	}

	var opts dagger.ContainerWithExecOpts
	if svc.Privileged {
		opts.InsecureRootCapabilities = true
	}

	ctr = ctr.WithExec(svc.Command, opts)

	return &Service{
		Name:      svc.Name,
		Container: ctr,
	}, nil
}
