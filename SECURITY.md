<!--
Copyright (c) 1898 & Co.
SPDX-License-Identifier: Apache-2.0
-->

# Security Policy

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in this project, please report it responsibly.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to [security@1898andco.com](mailto:security@1898andco.com).

Include the following information in your report:

- Type of vulnerability (e.g., authentication bypass, injection, etc.)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the vulnerability and how it might be exploited

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Communication**: We will keep you informed of progress toward a fix
- **Disclosure**: We will coordinate with you on the timing of public disclosure
- **Credit**: We will credit you in the security advisory (unless you prefer to remain anonymous)

### Safe Harbor

We consider security research conducted in accordance with this policy to be:

- Authorized and lawful
- Helpful to the security of our users
- Conducted in good faith

We will not pursue legal action against researchers who follow this policy.

## Security Best Practices

When using this SDK:

1. **Protect API Keys**: Never commit API keys to version control. Use environment variables or secure secret management
2. **Use HTTPS**: The SDK defaults to HTTPS. Do not override this with HTTP in production
3. **Rotate Credentials**: Regularly rotate API keys according to your organization's security policy
4. **Monitor Access**: Enable audit logging in Armis to monitor API usage
5. **Least Privilege**: Use API keys with the minimum required permissions
