function Gmail_Trigger(){   
    const TRIGGER_LABEL = 'MessagesToHandle'
    const HANDLED_LABEL = 'HandledMessages'                          
    const SUBSCRIBE = '<your-email>+subscribe@gmail.com'
    const UNSUBSCRIBE = '<your-email>+unsubscribe@gmail.com'
    const API_USER = 'admin'
    const API_PASSWORD = '<your-admin-password>'
    const DOMAIN = 'pupdate.today'
  var messages = GmailApp.get
    //-----------------------------------------------------------//
   
    var label = GmailApp.getUserLabelByName(TRIGGER_LABEL);  
    var handledLabel = GmailApp.getUserLabelByName(HANDLED_LABEL);  
    if (label == null){
    GmailApp.createLabel(triggerLabel);
    } else{
  
    var url = "https://" + DOMAIN + "/createToken";
    var params = {
              method: "POST",
              contentType: 'application/json',
              headers: {
                'Authorization': 'Basic ' + Utilities.base64Encode(API_USER + ':' + API_PASSWORD)
              }
          };
  
    console.log("-----Fetch Token-----")
    const TOKEN = UrlFetchApp.fetch(url, params);
    console.log("-----TOKEN-----")
    console.log(TOKEN.getResponseCode())
    
    // Need the following regex to remove and replace the wrapped quotes
    const JWT = TOKEN.getContentText().replace(/^"|"$/g, '');
  
    var threads = label.getThreads();
    console.log("Starting Iteration of: %d ", threads.length)
    for(i=0;i<(threads.length); i++) {
      var messages = threads[i].getMessages()
      var recipients = messages[0].getTo().split(",");
      var sender = messages[0].getFrom().split("<")[1].split(">")[0];
      for(j=0;j<(recipients.length); j++) {
        if (recipients[j] == SUBSCRIBE){
          console.log("Subscribe: %s", sender)
          var url = "https://" + DOMAIN + "/subscriber";
          var data = {
            email: sender
          };
          var options = {
              method: "post",
              contentType: 'application/json',
              "payload": JSON.stringify(data),
              headers: {
                'Authorization': 'Bearer ' + JWT
              }
          };

          var response = UrlFetchApp.fetch(url, options);
          console.log("RESP: %s", response.getResponseCode());
        }

        else if (recipients[j] == UNSUBSCRIBE) {
          console.log("Unsubscribe: %s", sender)
          var url = "https://" + DOMAIN + "/subscriber/" + sender
          var params = {
              method: "DELETE",
               contentType: 'application/json',
               headers: {
                'Authorization': 'Bearer ' + JWT
              }
          };
          console.log("URL: "+ url)
          const responseString = UrlFetchApp.fetch(url, params);
          console.log("Status %s", responseString.getResponseCode());
        }
      }
      // Mark that thread is handled
      threads[i].removeLabel(label)
      threads[i].addLabel(handledLabel)
      }
    }
  }