# driftwatch

> Detect configuration drift between running services and their declared infrastructure-as-code definitions.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your IaC definitions and a running environment to detect drift:

```bash
# Scan a Terraform state against live AWS resources
driftwatch scan --provider aws --iac ./terraform --region us-east-1

# Output results as JSON
driftwatch scan --provider aws --iac ./terraform --format json

# Watch for drift continuously (every 5 minutes)
driftwatch watch --provider aws --iac ./terraform --interval 5m
```

Example output:

```
[DRIFT] aws_instance.web-server
  Expected: instance_type = t3.micro
  Actual:   instance_type = t3.large

[OK]    aws_s3_bucket.assets
[DRIFT] aws_security_group.app — 2 rule(s) differ

Summary: 2 drifted, 1 clean
```

---

## Supported Providers

- AWS (Terraform)
- GCP (Terraform)
- Kubernetes (Helm / raw manifests)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

[MIT](LICENSE)