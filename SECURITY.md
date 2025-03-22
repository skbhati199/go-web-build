# Security Policy

## Supported Versions

The following versions of Go Web Build are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 3.0.x   | :white_check_mark: |
| 2.0.x   | :white_check_mark: |
| 1.0.x   | :x:                |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of Go Web Build seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Reporting Process

1. **Do NOT report security vulnerabilities through public GitHub issues.**

2. **Email us** at [skbhati199@gmail.com](mailto:skbhati199@gmail.com) with:
   - Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
   - Full paths of source file(s) related to the manifestation of the issue
   - Location of the affected source code (tag/branch/commit or direct URL)
   - Any special configuration required to reproduce the issue
   - Step-by-step instructions to reproduce the issue
   - Proof-of-concept or exploit code (if possible)
   - Impact of the issue, including how an attacker might exploit it

### Response Process

- You will receive an acknowledgment within 48 hours
- Our security team will investigate and provide updates every 5 business days
- If the issue is confirmed:
  - We will release a patch as soon as possible depending on complexity
  - We will publish a security advisory on GitHub
- If the issue is declined:
  - We will provide a detailed explanation of why

## Security Updates

Security updates will be released through our normal release process:

1. Patches will be released for all supported versions
2. Security advisories will be published on GitHub
3. Release notes will include security-related fixes

## Best Practices

When using Go Web Build in your projects:

1. Always use the latest supported version
2. Regularly check for security advisories
3. Follow our security guidelines in the documentation
4. Keep all dependencies up to date
5. Use security linting tools
6. Enable all recommended security features

## Security-Related Configuration

For secure deployment, ensure:

1. Proper CORS configuration
2. Secure cookie settings
3. HTTPS enforcement
4. Content Security Policy (CSP) headers
5. Rate limiting
6. Input validation

## Acknowledgments

We would like to thank the following security researchers who have helped improve Go Web Build's security:

- List will be updated as contributors report issues

## Contact

For any security-related questions, contact us at:
- Email: skbhati199@gmail.com
- PGP Key: [Link to PGP key]
