# Tag Selection

By default, watchtower will update your container to the latest available image that matches the current tag of your running container. For example, if you are running `myimage:latest`, watchtower will check for a newer version of `myimage:latest`.

However, you might want watchtower to automatically switch to a newer versioned tag (e.g., from `v1.0.1` to `v1.1.0`) if one is available in the registry. This can be achieved using **Tag Selection Strategies**.

## Tag Selection Labels

You can configure tag selection on a per-container basis using the following labels:

- `com.centurylinklabs.watchtower.tag-strategy`: The strategy to use for selecting the best tag. Currently, only `semver` is supported.
- `com.centurylinklabs.watchtower.tag-filter`: An optional regular expression to filter the tags before applying the strategy.

## SemVer Strategy

The `semver` strategy allows watchtower to find the highest versioned tag available in the registry based on [Semantic Versioning](https://semver.org/).

### How it works

1. Watchtower fetches all available tags for the container's image from the registry.
2. If a `tag-filter` is provided, it filters the tags using the regular expression.
3. It parses the remaining tags as semantic versions.
4. It selects the highest version and updates the container to use that tag.

### Coercion and Prefix Support

The `semver` strategy is flexible but has specific rules for prefixes:

- **Supported Prefixes**: Only the letter `v` (case-insensitive) is supported as a prefix (e.g., `v1.0.0`).
- **Unsupported Prefixes**: Custom prefixes like `demo1.0.0` or `app-2.0` will result in a parsing error and the tag will be skipped.
- **Auto-completion**: 2-part versions are automatically completed (e.g., `1.0` becomes `1.0.0`).

## Performance & Best Practices

Using tag selection involves an $O(N)$ operation where $N$ is the number of tags in your registry repository. For large repositories with thousands of tags, this can cause performance hits or registry rate-limiting.

### Infrastructure Optimization

To keep scans fast and efficient:

- **Repository Sharding**: Separate your "Development" and "Production" images into different repositories. This keeps the number of tags per repository small.
- **Tag Pruning**: Periodically delete old or unused tags from your registry.
- **Specific Filtering**: Use a narrow `tag-filter` regex. While watchtower still has to list all tags, a specific regex reduces the amount of strings that need to be parsed by the SemVer engine.

### GitOps for High-Frequency Releases

If your project has a high-velocity release cycle or requires a strict audit trail for every deployment, consider using a GitOps tool like [Watcher](https://github.com/Sithukyaw666/watcher).

While Watchtower is excellent for background image updates, Watcher is designed specifically for Docker Compose GitOps. It uses your Git repository as the "Source of Truth," ensuring that every deployment is deterministic and tracked in your version control history.

## Examples

#### Automatic Minor/Patch Updates

=== "docker-compose"

    ```yaml
    services:
      myapp:
        image: myorg/myapp:v1.0.0
        labels:
          com.centurylinklabs.watchtower.tag-strategy: semver
          com.centurylinklabs.watchtower.tag-filter: "^v1\\.\\d+\\.\\d+$"
    ```

#### Flexible Versioning (Coercion)

=== "docker-compose"

    ```yaml
    services:
      myapp:
        image: myorg/myapp:v0.1
        labels:
          com.centurylinklabs.watchtower.tag-strategy: semver
          com.centurylinklabs.watchtower.tag-filter: "^v\\d+\\.\\d+$"
    ```
