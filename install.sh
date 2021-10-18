#!/usr/bin/env bash

GO_HOOKS_DIR_NAME=".go-hooks"
GO_HOOKS_DIR=$HOME/$GO_HOOKS_DIR_NAME
GO_HOOKS_CONFIG=$GO_HOOKS_DIR/config.yaml

GO_HOOKS_EXECUTABLE="gohooks"

GO_HOOKS_EXECUTABLE_PATH="$HOME/.local/bin/$GO_HOOKS_EXECUTABLE"

RELEASE_URL="https://github.com/UsingCoding/gohooks/releases/download/v0.1.2/gohooks"

GIT_HOOKS=("commit-msg" "pre-push")

# Colors

RED='\033[0;31m'
YELLOW='\033[1;33m'
LIGHT_CYAN='\033[1;36m'

CC='\033[0m' # Clear Color

gitHookTemplate() {
    HOOK_NAME=$1

    cat <<-EOF
#!/usr/bin/env bash

$GO_HOOKS_EXECUTABLE hook $HOOK_NAME "\$@"
EOF
}

configTemplate() {
    cat <<-EOF
protectedReposRegExps:
# NOTE: regexp must fully match remote address
# Regexp to check will be chosen by remote name
# If remote url satisfy any of regexp hook applied
#    origin:
# Apply for repos from github
#        - .*github.com.*
# Apply only for ssh remotes
#        - git@+*
#    source:
# Apply for repos with specific names
#        - .*SpecificName.*
EOF
}

# Installing gohooks binary at $PATH
curl -Ls $RELEASE_URL -o "$GO_HOOKS_EXECUTABLE_PATH" || exit 1
chmod +x "$GO_HOOKS_EXECUTABLE_PATH"
echo -e "${LIGHT_CYAN}$GO_HOOKS_EXECUTABLE installed at $GO_HOOKS_EXECUTABLE_PATH${CC}"

if [ ! -d "$GO_HOOKS_DIR" ]; then
    echo "Creating $GO_HOOKS_DIR"
    mkdir "$GO_HOOKS_DIR" || exit 1
fi

echo -e "${LIGHT_CYAN}Note:${CC} This will rewrite git hooks files"

for GIT_HOOK in "${GIT_HOOKS[@]}"; do
    echo "..Set up $GIT_HOOK"

    echo "....Creating $GIT_HOOK"

    HOOK_PATH="$GO_HOOKS_DIR/$GIT_HOOK"

    TEMPLATE=$(gitHookTemplate "$GIT_HOOK")

    echo "$TEMPLATE" > "$HOOK_PATH"

    echo "....Setup $GIT_HOOK as executable"

    chmod u+x "$HOOK_PATH"
done

if [ ! -f "$GO_HOOKS_CONFIG" ]; then
    echo -e "${LIGHT_CYAN}Config not found:${CC} creating new one"

    configTemplate > "$GO_HOOKS_CONFIG"

    if [[ -z "$EDITOR" ]]; then
       echo -e "${RED}You are a bad linux user: your env \$EDITOR is empty:${CC} edit $GO_HOOKS_CONFIG by your own way"
   else
       $EDITOR "$GO_HOOKS_CONFIG"
    fi

fi

echo -e "${LIGHT_CYAN}go-hooks as global git hooks installed${CC}"
git config --global core.hooksPath "$GO_HOOKS_DIR"