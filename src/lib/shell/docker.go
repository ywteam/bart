package shell

type Docker struct {
}

func (d Docker) Cli(args ...string) *ShellCommand {
	if !Features["docker"] {
		panic("Docker is not available on this system.")
	}
	return NewCommand("docker", args...)
}
func (d Docker) Run(image string, opts DockerRunOptions) *ShellCommand {
	args := []string{"run"}
	if opts.ContainerName != "" {
		args = append(args, "--name", opts.ContainerName)
	}
	if opts.Volumes != nil {
		for _, v := range opts.Volumes {
			args = append(args, "-v", v[0]+":"+v[1])
		}
	}
	if opts.Network != nil {
		for _, n := range opts.Network {
			args = append(args, "--network", n)
		}
	}
	if opts.Ports != nil {
		for _, p := range opts.Ports {
			args = append(args, "-p", p[0]+":"+p[1])
		}
	}
	if opts.Env != nil {
		for k, v := range opts.Env {
			args = append(args, "-e", k+"="+v)
		}
	}
	if opts.Detach {
		args = append(args, "-d")
	}
	if opts.Interactive {
		args = append(args, "-i")
	}
	if opts.Tty {
		args = append(args, "-t")
	}
	if opts.Remove {
		args = append(args, "--rm")
	}
	if opts.Privileged {
		args = append(args, "--privileged")
	}
	if opts.Workdir != "" {
		args = append(args, "-w", opts.Workdir)
	}
	if opts.User != "" {
		args = append(args, "-u", opts.User)
	}
	if opts.Entrypoint != "" {
		args = append(args, "--entrypoint", opts.Entrypoint)
	}
	args = append(args, image)
	if opts.Cmd != "" {
		args = append(args, opts.Cmd)
	}
	if opts.CmdArgs != nil {
		args = append(args, opts.CmdArgs...)
	}
	return d.Cli(args...)
}
func (d Docker) Exec(container string, cmd string, args ...string) *ShellCommand {
	return d.Cli(append([]string{"exec", container, cmd}, args...)...)
}

type DockerRunOptions struct {
	ContainerName string
	Image         string
	Volumes       [][]string
	Network       []string
	Ports         [][]string
	// Env           []string
	Env         map[string]string
	Args        []string
	Detach      bool
	Interactive bool
	Tty         bool
	Remove      bool
	Privileged  bool
	Workdir     string
	User        string
	Entrypoint  string
	Cmd         string
	CmdArgs     []string
}
