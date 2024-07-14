# Project Title

Tenant Management Module for SaaS

## Description

API Backend Tenant Management Module for SaaS

## Getting Started

### Dependencies

- Disarankan menggunakan WSL2 dengan distro Ubuntu 22.04
- [Go](https://go.dev/doc/install) (Gunakan versi terbaru)

### Installing

1. Clone repository:

   ```bash
   git clone git@github.com:Marcellinom/tenant-management-saas.git
  
   # Untuk pengembangan
   git checkout development
   ```

### Executing program

1. Copy .env.example ke .env
   ```bash
   cp .env.example .env
   ```
2. Set up environment variable `.env`
3. Jalankan server.
   ```bash
   go run .
   ```
## License

This project is licensed under the MIT License - see the [LICENSE file](./LICENSE) for details

## Acknowledgments
- [Terraform](https://www.terraform.io/)
- [Base Go DPTSI](https://bitbucket.org/dptsi/base-go)
- many others...
