#!/bin/bash

# Script Name: git-curl.sh
# Description: A script to demonstrate handling command-line arguments.
set -e
VERSION="1.0.0"

# Function to securely read access token
read_access_token() {
    echo -n "Enter your GitHub access token: "
    read -s ACCESS_TOKEN
    echo
}

# Function to handle setup command
setup_command() {
    if [[ $# -ne 3 ]]; then
        echo "Usage: $0 setup <user name> <folder path for code> <git repo url>"
        exit 1
    fi

    USER_NAME="$1"
    FOLDER_PATH="$2"
    GIT_REPO_URL="$3"

    # Create .netrc file
    read_access_token
    echo "machine github.com login $USER_NAME password $ACCESS_TOKEN" > ~/.netrc
    chmod 600 ~/.netrc
    echo ".netrc file created."

    # Create .config file
    mkdir -p "$FOLDER_PATH"
    echo "repo_url=$GIT_REPO_URL" > "$FOLDER_PATH/.config"
    echo "branch=master" >> "$FOLDER_PATH/.config"
    echo ".config file created in $FOLDER_PATH."
}

# Function to handle get command
get_command() {
    if [[ $# -ne 2 ]]; then
        echo "Usage: $0 get <branch|tag|commit> <name>"
        exit 1
    fi

    TYPE="$1"
    NAME="$2"

    if [[ ! -f .config ]]; then
        echo ".config file not found in the current directory."
        exit 1
    fi

    # Read repo_url and branch from .config
    source .config
    if [[ -z "$repo_url" ]]; then
        echo "repo_url not found in .config file."
        exit 1
    fi

    # Determine the URL based on type
    case "$TYPE" in
        branch)
            DOWNLOAD_URL="$repo_url/archive/refs/heads/$NAME.tar.gz"
            ;;
        tag)
            DOWNLOAD_URL="$repo_url/archive/refs/tags/$NAME.tar.gz"
            ;;
        commit)
            DOWNLOAD_URL="$repo_url/archive/$NAME.tar.gz"
            ;;
        *)
            echo "Invalid type: $TYPE. Use branch, tag, or commit."
            exit 1
            ;;
    esac

    # Update branch in .config if type is branch
    if [[ "$TYPE" == "branch" ]]; then
        sed -i.bak "s/^branch=.*/branch=$NAME/" .config
    fi

    # Download and extract the tar.gz
    echo "Downloading $DOWNLOAD_URL..."
    curl -L -n -o code.tar.gz "$DOWNLOAD_URL"
    if [[ $? -ne 0 ]]; then
        echo "Failed to download $DOWNLOAD_URL."
        exit 1
    fi

    # Clean up source folder and extract
    rm -rf source
    mkdir source
    tar -xzf code.tar.gz -C source --strip-components=1
    rm code.tar.gz
    echo "Code extracted to source folder."
}

show_help() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  setup <user name> <folder path for code> <git repo url>"
    echo "      Set up the environment by creating necessary configuration files."
    echo ""
    echo "  get <branch|tag|commit> <name>"
    echo "      get the specified branch, tag, or commit from the repository."
    echo ""
    echo "  version    Show the script version and exit"
    echo ""
    echo "Options:"
    echo "  --help       Show this help message and exit"
    echo ""
}

# Function to display version
show_version() {
    echo "$0 version $VERSION"
}

# Parse command-line arguments
if [[ $# -eq 0 ]]; then
    echo "No arguments provided. Use --help for usage information."
    exit 1
fi

COMMAND="$1"
shift

case "$COMMAND" in
    setup)
        setup_command "$@"
        ;;
    get)
        get_command "$@"
        ;;
    --help)
        show_help
        ;;
    version)
        show_version
        ;;
    *)
        echo "Unknown command: $COMMAND"
        echo "Use --help for usage information."
        exit 1
        ;;
esac