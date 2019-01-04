# rgv

_**Reddit Graph Visualizer:** A tool for visualizing the relationships amongst
the users of a subreddit._

[![Github: Tag][tag-img]][tag]
[![Travis: Build][travis-img]][travis]
[![Codecov: Coverage][codecov-img]][codecov]
[![Go Report Card][grp-img]][grp]

```
This is a work-in-progress (under active development).
```

## Usage

We're live! Check it out at https://rgv.stevenxie.me.

_TODO: Add more information about features._

## Roadmap

- [x] Develop initial realtime-only version.
- [ ] Make visualizer more informative and usable.
- [ ] Improve / stabilize backend infrastructure.
- [ ] Add database-connected history-tracking capabilities.

<br />

## Deployment

`rgv` is configured to build as a set of two Docker images (an entrypoint and
an API server), and to deploy on Kubernetes.

To deploy `rgv`, create the Kubernetes resources defined in `deployment/`:

```bash
cat deployment/* | kubectl create -f -
```

### Exposure Options

#### Load Balancer:

> This is the easier way to make `rgv` publicly accessible; however, this
> requires the use of a load balancer, which may be cost-prohibitive.

If you intend to run `rgv` behind a load balancer, go ahead and change the
`rgv-entrypoint` service to type `LoadBalancer`:

```bash
kubectl patch service rgv-entrypoint -p "
spec:
  type: LoadBalancer
"
```

#### Ingress Controller:

> This method requires for you to have preconfigured an Ingress controller, like
> [Traefik](https://docs.traefik.io/user-guide/kubernetes/), which will route
> traffic from a publicly exposed node to the `rgv-entrypoint` service.

If you intend to run `rgv` behind an Ingress controller (what I do), go ahead
and create an Ingress resource:

```bash
kubectl create -f - <<EOF
apiVersion: extensions/v1beta1
kind: Ingress

metadata:
  name: rgv-entrypoint

spec:
  rules:
    - host: "$YOUR_HOST"  # e.g. rgv.example.com or an external IP
      http:
        paths:
          - backend:
              serviceName: rgv-entrypoint
              servicePort: http
EOF
```

This Ingress resource is defined to route traffic from `$YOUR_HOST` to the
`http` port on the `rgv-entrypoint` service (port 80).

### Automatic HTTPS

I personally do not enable HTTPS on `rgv` directly, since I have it sitting
behind a [Traefik](https://traefik.io) edge router that takes care of SSL
certificate management for me.

However, in the event that `rgv-entrypoint` is directly exposed to an external
network, it does indeed support automatic SSL certificate management thanks to
[Caddy](https://caddyserver.com) and [Let's Encrypt](https://letsencrypt.org).
Simply create a `ConfigMap` resource to set it up:

```bash
kubectl create -f - <<EOF
apiVersion: v1
kind: ConfigMap

metadata:
  name: rgv-entrypoint-config

data:
  HOST: "$YOUR_HOST"        # e.g. rgv.example.com or an external IP
  TLS_STRATEGY: "email"     # one of 'email' or 'dns'
  TLS_EMAIL: "$YOUR_EMAIL"  # e.g. email@example.com
EOF
```

This will make `rgv-entrypoint` attempt to generate an SSL certificate using
`$TLS_EMAIL` for `$YOUR_HOST` upon startup.

However, there are some extra steps you will have to take in order to make
those certificates persist; see
[`deployment/optionals/`](https://github.com/stevenxie/rgv/tree/master/deployment/optionals)
for details.

[tag]: https://github.com/stevenxie/rgv/releases
[tag-img]: https://img.shields.io/github/tag/stevenxie/rgv.svg
[travis]: https://travis-ci.com/stevenxie/rgv
[travis-img]: https://travis-ci.com/stevenxie/rgv.svg?branch=master
[codecov]: https://codecov.io/gh/stevenxie/rgv
[codecov-img]: https://codecov.io/gh/stevenxie/rgv/branch/master/graph/badge.svg
[grp]: https://goreportcard.com/report/github.com/stevenxie/rgv
[grp-img]: https://goreportcard.com/badge/github.com/stevenxie/rgv
