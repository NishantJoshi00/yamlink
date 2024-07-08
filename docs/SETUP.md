Here's a rewritten version of your setup that's easier for others to follow:

**Step 1: Install Yamlink**

- Download and install the `yamlink` executable in a directory like `/usr/bin/`.
- Make sure you have the necessary permissions to run the executable.

**Step 2: Create a Systemd Unit File**

- Create a new file called `yamlink.service` in the `/etc/systemd/system/` directory with the following contents:

```ini
[Unit]
Description=Yamlink Service
After=network.target

[Service]
Type=simple
Environment=CONFIG_FILE="/etc/yamlink/config.yaml"
ExecStart=/usr/bin/yamlink
Restart=always

[Install]
WantedBy=multi-user.target
```

- Save the file and reload the systemd daemon by running `sudo systemctl daemon-reload`.
- Enable the Yamlink service to start automatically on boot by running `sudo systemctl enable yamlink`.

**Step 3: Create a Configuration File**

- Create a new file called `config.yaml` in the `/etc/yamlink/` directory with the following contents:

```yaml
host: 127.0.0.1
port: 80
map_file: /etc/yamlink/mapping.yaml
refresh_interval: 5
```

- Save the file and make sure it has the correct permissions.

**Step 4: Create a Mapping File**

- Create a new file called `mapping.yaml` in the `/etc/yamlink/` directory with the following contents:

```yaml
github:
  profile: https://github.com/NishantJoshi00
```

- Save the file and make sure it has the correct permissions.

**Step 5: Add Server to Hosts File**

- Edit your system's hosts file (usually `/etc/hosts`) and add the following line:

```bash
127.0.0.1 s
```

**Step 6: Use Yamlink**

- Once the setup is complete, you can access Yamlink by typing `s/<service name>/` in your browser.
- For example, if you want to access your GitHub profile, you would type `s/github/profile/`.

Tips:

- Make sure to replace `<service name>` with the actual domain name you specified in the hosts file (e.g. `s`).
- You can update the mapping file and Yamlink will automatically pick up the changes.
- You can customize the configuration file to suit your needs.

By following these steps, others should be able to set up their own Yamlink service and start using it productively!
