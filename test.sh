#!/bin/bash

## Settings
RP_REPO_URL="https://github.com/rocket-pool/rocketpool.git"
RP_REPO_BRANCH="v1.2"
HARDHAT_PORT=8545


## Dependencies
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
source ~/.bashrc
nvm install v20
nvm use v20


##
# Helpers
##


# Clean up
cleanup() {

    # Remove RP repo
    if [ -d "$RP_TMP_PATH" ]; then
        rm -rf "$RP_TMP_PATH"
    fi

    # Kill ganache instance
    if [ -n "$HARDHAT_PID" ] && ps -p "$HARDHAT_PID" > /dev/null; then
        kill -9 "$HARDHAT_PID"
    fi

}

# Clone the contract repos
clone_repos() {
    RP_TMP_PATH="$(mktemp -d)"
    
    RP_PATH="$RP_TMP_PATH/rocketpool"
    git clone "$RP_REPO_URL" -b "$RP_REPO_BRANCH" "$RP_PATH"

    MULTICALL_PATH="$RP_TMP_PATH/multicall"
    git clone "$MULTICALL_REPO_URL" "$MULTICALL_PATH"

    BALANCE_BATCHER_PATH="$RP_TMP_PATH/eth-balance-checker"
    git clone "$BALANCE_BATCHER_REPO_URL" -b "$BALANCE_BATCHER_REPO_BRANCH" "$BALANCE_BATCHER_PATH"
}

# Install Rocket Pool dependencies
install_rp_deps() {
    cd "$RP_PATH"
    npm install
    cd - > /dev/null
}

# Start the hardhat EVM and server
start_hardhat() {
    cd "$RP_PATH"
    npx hardhat node --port $HARDHAT_PORT > /dev/null &
    HARDHAT_PID=$!
    cd - > /dev/null
}

# Deploy Rocket Pool contracts
deploy_rp() {
    cd "$RP_PATH"
    npx hardhat run --network localhost scripts/deploy.js
    cd - > /dev/null
}

# Install dependencies for the test lib
install_test_deps() {
    cd tests/hardhat
    npm install
    cd - > /dev/null
}

# Deploy other test contracts
deploy_other() {
    cd tests/hardhat
    npx hardhat run --network localhost scripts/deploy.js
    cd - > /dev/null
}

# Run tests
run_tests() {
    go clean -testcache && go test -p 1 --tags=testing ./...
}


##
# Run
##

# Clean up before exiting
trap cleanup EXIT

# Clone RP repo
echo ""
echo "Cloning main Rocket Pool repository..."
echo ""
clone_rp

# Install RP deps
echo ""
echo "Installing Rocket Pool dependencies..."
echo ""
install_rp_deps

# Start hardhat
echo ""
echo "Starting the hardhat server..."
echo ""
start_hardhat

# Deploy RP contracts
echo ""
echo "Deploying Rocket Pool contracts..."
echo ""
deploy_rp

# Install RP deps
echo ""
echo "Installing other testing dependencies..."
echo ""
install_test_deps

# Deploy other contracts
echo ""
echo "Deploying other contracts..."
echo ""
deploy_other

# Run tests
echo ""
echo "Running tests..."
echo ""
run_tests
