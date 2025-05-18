GITHUB_API_HOST="api.github.com"

# Helper function to parse JSON and extract filenames
# Usage: parse_filenames <json> <jq_filter>
parse_filenames() {
    local jq_filter="$1"
    if command -v jq >/dev/null 2>&1; then
        jq -r "$jq_filter"
    else
        python3 -c "
import sys, json
data = json.load(sys.stdin)
def extract(obj, path):
    for key in path.split('.'):
        if isinstance(obj, list):
            obj = [item.get(key, {}) for item in obj]
        else:
            obj = obj.get(key, {})
    if isinstance(obj, list):
        for item in obj:
            if isinstance(item, str):
                print(item)
    elif isinstance(obj, str):
        print(obj)
if '$jq_filter' == '.[].sha':
    for item in data:
        print(item.get('sha', ''))
elif '$jq_filter' == '.files[].filename':
    for item in data.get('files', []):
        print(item.get('filename', ''))
elif '$jq_filter' == '.commits[].sha':
    for item in data.get('commits', []):
        print(item.get('sha', ''))
"
    fi
}

#   github_modified_files_by_date <owner> <repo> <from_date> <to_date>
#   github_modified_files_by_sha <owner> <repo> <from_sha> <to_sha>
#
# Date format: YYYY-MM-DDTHH:MM:SSZ (ISO 8601, e.g., 2024-06-01T00:00:00Z)
# .netrc file must be configured for api.github.com:
# machine api.github.com
#   login <your_github_username>
#   password <your_github_token>


github_modified_files_by_date() {
    local CFG="$1"
    local from_date="$2" # Format: YYYY-MM-DDTHH:MM:SSZ
    local to_date="$3"   # Format: YYYY-MM-DDTHH:MM:SSZ

    if [[ -z "$from_date" || -z "$to_date" ]]; then
        echo "function github_modified_files_by_date params missing: <config_file> <from_date> <to_date>"
        exit 1
    fi

    if [[ ! -f "$CFG" ]]; then
        echo "$CFG file not found in the current directory."
        exit 1
    fi
    source "$CFG"
    # .config has branch, we need to use function parameter branch
    local branch="$4"

    # extract the github token from the .netrc file
    local github_token=$(awk -v machine="github.com" '$1 == "machine" && $2 == machine {print $6}' ~/.netrc)
    # echo "GitHub token: $github_token"

    curl -s -H "Authorization: Bearer $github_token" \
        "https://$GITHUB_API_HOST/repos/$repo_owner/$repo_name/commits?since=$from_date&until=$to_date&sha=$branch&per_page=100" |
        parse_filenames '.[].sha' |
        while read sha; do
            # echo "Processing commit: $sha"
            curl -s -H "Authorization: Bearer $github_token" \
                "https://$GITHUB_API_HOST/repos/$repo_owner/$repo_name/commits/$sha" |
                parse_filenames '.files[].filename'
        done | sort -u
}

github_modified_files_by_sha() {
    local CFG="$1"
    local from_sha="$2"
    local to_sha="$3"

    if [[ -z "$from_sha" || -z "$to_sha" || -z "$branch" ]]; then
        echo "function github_modified_files_by_sha params missing: <config_file> <from_sha> <to_sha> <branch>"
        exit 1
    fi

    if [[ ! -f "$CFG" ]]; then
        echo "$CFG file not found in the current directory."
        exit 1
    fi
    source "$CFG"
    # .config has branch, we need to use function parameter branch
    local branch="$4"

    # extract the github token from the .netrc file
    local github_token=$(awk -v machine="github.com" '$1 == "machine" && $2 == machine {print $6}' ~/.netrc)
    # echo "GitHub token: $github_token"
    curl -s -H "Authorization: Bearer $github_token" \
        "https://$GITHUB_API_HOST/repos/$repo_owner/$repo_name/compare/$from_sha...$to_sha" |
        parse_filenames '.commits[].sha' |
        while read sha; do
            curl -s -H "Authorization: Bearer $github_token" \
                "https://$GITHUB_API_HOST/repos/$repo_owner/$repo_name/commits/$sha" |
                parse_filenames '.files[].filename'
        done | sort -u
}
