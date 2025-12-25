---
title: "Configuration"
description: "Configure your site settings and options"
order: 3
---

Learn how to configure your site using the `site.yml` configuration file.

## Configuration File

Create a `site.yml` file in your project root:

```yaml
title: "My Site"
description: "A description of my site"
base_url: "https://example.com"
```

## Available Options

### Site Information

| Option | Description | Required |
|--------|-------------|----------|
| `title` | Your site's title | Yes |
| `description` | Site description for SEO | Yes |
| `base_url` | The canonical URL of your site | Yes |

### Authentication (Optional)

For emoji reactions, you'll need to configure Google OAuth:

```yaml
# Set as environment variables
# GOOGLE_CLIENT_ID=your-client-id
# GOOGLE_CLIENT_SECRET=your-client-secret
```

> [!WARNING]
> Never commit OAuth credentials to version control. Use environment variables instead.

## Environment Variables

The following environment variables are supported:

- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `PORT` - Server port (default: 3000)

## Validation

Run the following to validate your configuration:

```bash
./site build --dry-run
```

This will check your configuration without generating any files.

## Next Steps

With your site configured, learn how to create content in the next section.
