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

### Coercion

The `semver` strategy is flexible and can handle versions that are not strictly SemVer-compliant by "coercing" them. For example:

- `v1.0` is treated as `1.0.0`
- `2.1` is treated as `2.1.0`
- `v1` is treated as `1.0.0`

### Examples

#### Automatic Minor/Patch Updates

If you want your container to automatically update to the latest `v1.x.x` version but never switch to `v2.x.x`, you can use a filter.

=== "docker-compose"

    ```yaml
    services:
      myapp:
        image: myorg/myapp:v1.0.0
        labels:
          com.centurylinklabs.watchtower.tag-strategy: semver
          com.centurylinklabs.watchtower.tag-filter: "^v1\\.\\d+\\.\\d+$"
    ```

=== "docker run"

    ```bash
    docker run -d \
      --label com.centurylinklabs.watchtower.tag-strategy=semver \
      --label com.centurylinklabs.watchtower.tag-filter="^v1\.\d+\.\d+$" \
      myorg/myapp:v1.0.0
    ```

In this example, if the registry has `v1.0.1`, `v1.1.0`, and `v2.0.0`, watchtower will identify `v1.1.0` as the highest semantic version and update your container.

#### Flexible Versioning

If you use tags like `v0.1` and `v1.0`:

=== "docker-compose"

    ```yaml
    services:
      myapp:
        image: myorg/myapp:v0.1
        labels:
          com.centurylinklabs.watchtower.tag-strategy: semver
          com.centurylinklabs.watchtower.tag-filter: "^v\\d+\\.\\d+$"
    ```

Watchtower will correctly identify `v1.0` as higher than `v0.1` even though they are missing the patch number.
