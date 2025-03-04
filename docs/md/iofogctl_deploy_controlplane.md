## iofogctl deploy controlplane

Deploy a Control Plane

### Synopsis

Deploy a Control Plane.

```
iofogctl deploy controlplane [flags]
```

### Examples

```
iofogctl deploy controlplane -f controlplane.yaml
```

### Options

```
  -f, --file string   YAML file containing resource definitions for Control Plane
  -h, --help          help for controlplane
```

### Options inherited from parent commands

```
      --config string      CLI configuration file (default is ~/.iofog/config.yaml)
      --http-verbose       Toggle for displaying verbose output of API client
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of iofogctl
```

### SEE ALSO

* [iofogctl deploy](iofogctl_deploy.md)	 - Deploy ioFog platform or components on existing infrastructure

###### Auto generated by spf13/cobra on 22-Oct-2019
