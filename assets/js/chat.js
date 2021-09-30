// Classes
const ChatUI = new Chat_UI(),
      ChatAPI = new Chat_API()

window.addEventListener('DOMContentLoaded', (event) => {
    const chatMessagesDOM = document.getElementById('chat-body')
    let currentMessages, oldMessages, difference

    for(var i = 0; i < 100; i++) {
        (function(index) {
            setTimeout(function() { 
                ChatAPI.getMessages().then(response => {
                    // oldMessages = currentMessages
                    // currentMessages = response
                    // if(currentMessages.length != chatMessagesDOM.childElementCount) {
                    //     let difference = currentMessages.filter(elm => !oldMessages.map(elm => JSON.stringify(elm)).includes(JSON.stringify(elm)));  
                    //     ChatUI.buildMessages(difference)
                    // }
                    ChatUI.buildMessages(response)
                })
             }, index*2000);
        })(i);
    }    
    
})