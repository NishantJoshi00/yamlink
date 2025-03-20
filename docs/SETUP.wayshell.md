**Setup Guide for `wayshell`**

**Prerequisites:**

* `wayshell` binary (see below for installation instructions)
* Shell configuration files (`.zshrc` or `config.fish`, depending on your shell)

**Installation:**

1. **Build the `wayshell` binary:**
```bash
go build ./cmd/wayshell
```
2. **Move the binary to a suitable location (optional):**
```bash
mv wayshell /usr/local/bin
```

**Configuration:**

1. **Create a configuration file:** Create a new file at `~/.wayshell.yaml` with the following content:
```yaml
gs: git status
```
Replace this example configuration with your own custom shortcuts!

**Shell Bindings:**

To set up shell bindings for `wayshell`, follow these steps:

* **For zsh users:**
Add the following line to your `.zshrc` file:
```bash
source scripts/init.zsh <path/to/wayshell> <path/to/wayshell-config>
```
Replace `<path/to/wayshell>` and `<path/to/wayshell-config>` with the actual paths on your system.
* **For fish users:**
Add the following line to your `config.fish` file:
```bash
source scripts/init.fish <path/to/wayshell> <path/to/wayshell-config>
```
Again, replace `<path/to/wayshell>` and `<path/to/wayshell-config>` with the actual paths on your system.

**Tips:**

* Make sure to adjust the path variables to match where you installed `wayshell` and its configuration file.
* You can customize the shortcuts in the `~/.wayshell.yaml` file to suit your workflow needs.

By following these steps, you should be able to set up `wayshell` and start using it to streamline your command-line experience!
