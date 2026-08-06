# Quickstart

## 1. Preview Your Docs

Run a local development server to see your changes in real-time:

```bash
npx @scalar/cli project preview
```

This starts a live preview at `http://localhost:7970` where every edit you make is instantly visible.

## 2. Customize Everything

Make it yours with themes, custom CSS, and MDX. Configure your documentation structure, navigation, and styling through `scalar.config.json`.

## 3. Publish Your Docs

First, authenticate with your Scalar account:

```bash
npx @scalar/cli auth login
```

Then publish your documentation:

```bash
npx @scalar/cli project publish --slug [dashboard-project-slug]
```

Your site will be available at `<your-slug>.apidocumentation.com`.

## Stuck?

Check whether your `scalar.config.json` is valid:

```bash
npx @scalar/cli project check-config
```

We're here to help:

- [Email support@scalar.com](mailto:support@scalar.com)
- [Chat with us on Discord](https://discord.gg/scalar)
- [Schedule a call](https://scalar.cal.com/scalar/chat-with-scalar)
