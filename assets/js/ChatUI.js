class Chat_UI {
    buildMessages(messages){
        const chatMessagesDOM = document.getElementById('chat-body')
        chatMessagesDOM.innerHTML = ''

        for (let i = 0; i < messages.length; i++) {
            const message = messages[i];
            
            const oneMessage = document.createElement('div')
            oneMessage.classList.add('message', 'info')
            // oneMessage.classList.add('message', 'my-message')
    
            oneMessage.innerHTML = ` 
                                        <img class="img-circle medium-image" src="/assets/images/${message.User}.png" alt="...">
                                        <div class="message-body">
                                            <div class="message-info">
                                                <h4> ${message.User} </h4>
                                                <h5> <i class="fa fa-clock-o"></i> ${message.Time} </h5>
                                            </div>
                                            <hr>
                                            <div class="message-text">
                                                ${message.Content}
                                            </div>
                                        </div>
                                        <br>`
            chatMessagesDOM.appendChild(oneMessage)   
        }
    }
}