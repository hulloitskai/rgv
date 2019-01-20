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

`rgv` is configured to build as a set of two Docker images (a frontend and an
API server), and to deploy on Kubernetes.

To deploy `rgv`, create the Kubernetes resources defined in `deployment/`:

```bash
cat deployment/* | kubectl create -f -
```

### Exposure Options

#### Load Balancer:

> This is the easier way to make `rgv` publicly accessible; however, this
> requires the use of a load balancer, which may be cost-prohibitive.

If you intend to run `rgv` behind a load balancer, go ahead and change the
`rgv-frontend` service to type `LoadBalancer`:

```bash
kubectl patch service rgv-frontend -p "
spec:
  type: LoadBalancer
"
```

#### Ingress Controller:

> This method requires for you to have preconfigured an Ingress controller, like
> [Traefik](https://docs.traefik.io/user-guide/kubernetes/), which will route
> traffic from a publicly exposed node to the `rgv-frontend` service.

If you intend to run `rgv` behind an Ingress controller (what I do), go ahead
and create an Ingress resource:

```bash
kubectl create -f - <<EOF
apiVersion: extensions/v1beta1
kind: Ingress

metadata:
  name: rgv-frontend

spec:
  rules:
    - host: "$YOUR_HOST"  # e.g. rgv.example.com or an external IP
      http:
        paths:
          - backend:
              serviceName: rgv-frontend
              servicePort: http
EOF
```

This Ingress resource is defined to route traffic from `$YOUR_HOST` to the
`http` port on the `rgv-frontend` service (port 80).

[tag]: https://github.com/stevenxie/rgv/releases
[tag-img]: https://img.shields.io/github/tag/stevenxie/rgv.svg
[travis]: https://travis-ci.com/stevenxie/rgv
[travis-img]: https://travis-ci.com/stevenxie/rgv.svg?branch=master
[codecov]: https://codecov.io/gh/stevenxie/rgv
[codecov-img]: https://codecov.io/gh/stevenxie/rgv/branch/master/graph/badge.svg
[grp]: https://goreportcard.com/report/github.com/stevenxie/rgv
[grp-img]: https://goreportcard.com/badge/github.com/stevenxie/rgv
