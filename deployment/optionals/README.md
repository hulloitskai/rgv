# Optionals

## Persistent ACME Certificates

In order for SSL certificates generated by `rgv` to persist between pod
instances, it needs a place to put them.

This is taken care of by a
[**Persistent Volume Claim**](https://kubernetes.io/docs/concepts/storage/persistent-volumes/):

```bash
kubectl create -f rgv-entrypoint-acme.pvc.yaml
```

You will need to patch the `rgv-entrypoint` service to mount
and use the claimed volume:

```bash
kubectl patch service rgv-entrypoint -p "
spec:
  template:
    spec:
      containers:
        - name: rgv-entrypoint
          volumes:
            - name: acme
              persistentVolumeClaim:
                claimName: rgv-entrypoint-acme
"
```