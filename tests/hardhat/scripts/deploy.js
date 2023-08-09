export const BalanceBatcher = artifacts.require('BalanceChecker.sol');
export const Multicall = artifacts.require('Multicall2.sol');

/*** Dependencies ********************/

const hre = require('hardhat');
const pako = require('pako');
const fs = require('fs');
const Web3 = require('web3');


/*** Utility Methods *****************/


// Compress / decompress ABIs
function compressABI(abi) {
    return Buffer.from(pako.deflate(JSON.stringify(abi))).toString('base64');
}
function decompressABI(abi) {
    return JSON.parse(pako.inflate(Buffer.from(abi, 'base64'), {to: 'string'}));
}

// Load ABI files and parse
function loadABI(abiFilePath) {
    return JSON.parse(fs.readFileSync(abiFilePath));
}


/*** Contracts ***********************/

// Multicall
const multicall = artifacts.require('Multicall2.sol');

// Balance Batcher
const balanceBatcher = artifacts.require('BalanceChecker.sol');


/*** Deployment **********************/


// Deploy the contracst for supporting Rocket Pool's test suite
export async function deployContracts() {
    // Set our web3 provider
    const network = hre.network;
    let $web3 = new Web3(network.provider);

    // Accounts
    let accounts = await $web3.eth.getAccounts(function(error, result) {
        if(error != null) {
            console.log(error);
            console.log("Error retrieving accounts.'");
        }
        return result;
    });

    console.log(`Using network: ${network.name}`);
    console.log(`Deploying from: ${accounts[1]}`)
    console.log('\n');

    // Deploy Multicall
    var multicallInstance = await multicall.new({from: accounts[1]});
    multicall.setAsDeployed(multicallInstance);
    const multicallAddress = (await multicall.deployed()).address;
    console.log('   Multicall Address');
    console.log('     ' + multicallAddress);
    
    // Deploy Balance Batcher
    var balanceBatcherInstance = await balanceBatcher.new({from: accounts[1]});
    balanceBatcher.setAsDeployed(balanceBatcherInstance);
    const balanceBatcherAddress = (await balanceBatcher.deployed()).address;
    console.log('   Balance Batcher Address');
    console.log('     ' + balanceBatcherAddress);

    // Log it
    console.log('\n');
    console.log('  Done!');
    console.log('\n');
};

// Run it
deployContracts().then(function() {
    process.exit(0);
});