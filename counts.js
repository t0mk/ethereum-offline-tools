function cab() { 
  var i = 0; 
  eth.accounts.forEach( function(e){
     console.log(e + " \tbalance: " + web3.fromWei(eth.getBalance(e),
        "ether") + " \tcount: " + eth.getTransactionCount(e) ); 
     i++; 
  });
}

cab();

