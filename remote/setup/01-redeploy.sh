#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Set the name of the new user to create.
USERNAME=omdb-api

# Prompt to enter a password for the PostgreSQL moviemaze user (rather than hard-coding
# a password in this script).
read -p "Enter DB_PASSWORD: " DB_PASSWORD
read -p "Enter mailtrap user: " MAILTRAP_USERNAME
read -p "Enter password for mailtrap user: " MAILTRAP_PASSWORD

# Force all output to be presented in en_US for the duration of this script. This avoids  
# any "setting locale failed" errors while this script is running, before we have 
# installed support for all locales. Do not change this setting!
export LC_ALL=en_US.UTF-8 

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

# Enable the "universe" repository.
add-apt-repository --yes universe

# Update all software packages.
apt update

if id "${USERNAME}" &>/dev/null; then
    echo "User ${USERNAME} already exists."
else
    # Add the new user and give them sudo privileges
    sudo useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

    # Force a password to be set for the new user the first time they log in.
    sudo passwd --delete "${USERNAME}"
    sudo chage --lastday 0 "${USERNAME}"

    # Copy the SSH keys from the template-user to the new user.
    sudo rsync --archive --chown=${USERNAME}:${USERNAME} /home/ew-user/.ssh /home/${USERNAME}
    
    echo "User ${USERNAME} created"
fi

# Install the goose db migration CLI tool. https://pressly.github.io/goose/installation/
curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

# Install PostgreSQL.
apt --yes install postgresql

# Set up the moviemaze DB and create a user account with the password entered earlier.
sudo -i -u postgres psql -c "CREATE DATABASE omdb_api"
sudo -i -u postgres psql -d omdb_api -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d omdb_api -c "CREATE ROLE omdb_api WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# Add a DSN and mailtrap to the system-wide environment variables in the /etc/environment file.
echo "OMDB_API_DB_DSN_PROD='postgres://omdb_api:${DB_PASSWORD}@localhost/omdb_api'" >> /etc/environment
echo "MAILTRAP_USERNAME='${MAILTRAP_USERNAME}'" >> /etc/environment
echo "MAILTRAP_PASSWORD='${MAILTRAP_PASSWORD}'" >> /etc/environment

# Install Caddy (see https://caddyserver.com/docs/install#debian-ubuntu-raspbian).
# apt install -y debian-keyring debian-archive-keyring apt-transport-https
# curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
# curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
# apt update
# apt --yes install caddy

# Upgrade all packages. Using the --force-confnew flag means that configuration 
# files will be replaced if newer ones are available.
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "Script complete! Rebooting..."
reboot