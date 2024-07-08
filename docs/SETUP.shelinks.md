**Setup Guide for `shelinks`**

**Prerequisites:**

* `shelinks` binary (see below for installation instructions)
* Shell configuration files (`.zshrc` or `config.fish`, depending on your shell)

**Installation:**

1. **Build the `shelinks` binary:**
```bash
go build ./cmd/shelinks
```
2. **Move the binary to a suitable location (optional):**
```bash
mv shelinks /usr/local/bin
```

**Configuration:**

1. **Create a configuration file:** Create a new file at `~/.shelinks.yaml` with the following content:
```yaml
gs: git status
```
Replace this example configuration with your own custom shortcuts!

**Shell Bindings:**

To set up shell bindings for `shelinks`, follow these steps:

* **For zsh users:**
Add the following line to your `.zshrc` file:
```bash
source scripts/init.zsh <path/to/shelinks> <path/to/shelinks-config>
```
Replace `<path/to/shelinks>` and `<path/to/shelinks-config>` with the actual paths on your system.
* **For fish users:**
Add the following line to your `config.fish` file:
```bash
source scripts/init.fish <path/to/shelinks> <path/to/shelinks-config>
```
Again, replace `<path/to/shelinks>` and `<path/to/shelinks-config>` with the actual paths on your system.

**Tips:**

* Make sure to adjust the path variables to match where you installed `shelinks` and its configuration file.
* You can customize the shortcuts in the `~/.shelinks.yaml` file to suit your workflow needs.

By following these steps, you should be able to set up `shelinks` and start using it to streamline your command-line experience!
