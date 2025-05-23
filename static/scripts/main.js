document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const section = document.querySelector('section');
  const mainPage = document.getElementById('mainPage');
  const usernameDisplay = document.getElementById('username');
  const logoutBtn = document.getElementById('logout');
  const body = document.getElementById('body');

  // Initialize the create post form
  section.innerHTML = `
    <form id="form" method="get">
      <div class="container">
        <h3>Create Post</h3>
        <div class="div-title">
          <label for="title">Title :</label>
          <input type="text" name="title" id="title" required>
        </div>
        <div class="div-description">
          <label for="description">description :</label>
          <textarea name="description" id="description" rows="4" required></textarea>
        </div>
        <div class="topic-options">
          <label><input type="checkbox" id="music" name="topic" value="Music"> Music</label>
          <label><input type="checkbox" id="sport" name="topic" value="Sport"> Sport</label>
          <label><input type="checkbox" id="gaming" name="topic" value="Gaming"> Gaming</label>
          <label><input type="checkbox" id="health" name="topic" value="Health"> Health</label>
          <label><input type="checkbox" id="general" name="topic" value="General"> General</label>
        </div>
        <div id="errorMsg" style="display:none; color:red; margin: 10px 10px;"></div>
        <button type="submit">Post</button>
      </div>
    </form>
  `;

  // Initialize the login/signup forms
  container.innerHTML = `
    <!-- Sign Up Form -->
    <div class="form-container sign-up-container">
      <form id="signupForm" action="/Signup" method="post">
        <h1>Create Account</h1>
        <input type="text" placeholder="Nickname" name="nickname" required />
        <input type="number" placeholder="Age" name="age" min="13" required />
        <select name="gender" required>
          <option value="" disabled selected>Gender</option>
          <option value="male">Male</option>
          <option value="female">Female</option>
          <option value="other">Prefer not to say</option>
        </select>
        <input type="text" placeholder="First Name" name="first_name" required />
        <input type="text" placeholder="Last Name" name="last_name" required />
        <input type="email" placeholder="Email" name="email" required />
        <input type="password" placeholder="Password" name="password" required />
        <input type="password" placeholder="Confirm Password" name="confirm_password" required />
        <button type="submit">Sign Up</button>
      </form>
    </div>

    <!-- Sign In Form -->
    <div class="form-container sign-in-container">
      <form id="loginForm">
        <input type="text" name="user" placeholder="Nickname" required />
        <input type="password" name="password" placeholder="Password" required />
        <button type="submit">Login</button>
      </form>
    </div>

    <!-- Overlay Panels -->
    <div class="overlay-container">
      <div class="overlay">
        <div class="overlay-panel overlay-left">
          <h1>Welcome Back!</h1>
          <p>If you already have an account, please sign in.</p>
          <button class="ghost" id="signIn">Sign In</button>
        </div>
        <div class="overlay-panel overlay-right">
          <h1>Hello, Friend!</h1>
          <p>Enter your details and start your journey with us.</p>
          <button class="ghost" id="signUp">Sign Up</button>
        </div>
      </div>
    </div>
  `;

  // Toggle between Sign Up and Sign In forms
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const loginForm = document.getElementById('loginForm');

  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    loginForm.style.display = 'none';
  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    loginForm.style.display = 'block';
  });

  // Check session on load
  checkSession();

  // Logout functionality
  logoutBtn.addEventListener('submit', async (e) => {
    e.preventDefault();
    await logout();
  });

  // Handle login form submission
  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    await handleLogin();
  });

  // Handle post creation form
  document.getElementById('form')?.addEventListener('submit', async function (e) {
    e.preventDefault();
    await createPost();
  });

  // Helper functions
  async function checkSession() {
    try {
      const response = await fetch('/check-session');
      const data = await response.json();

      if (data.loggedIn) {
        mainPage.style.display = 'block';
        currentUser = data.username; // Store current user's username
        usernameDisplay.textContent = `Welcome, ${data.username}!`;
        container.style.display = 'none';
        loadPosts();
        showMessagingUI();
      } else {
        container.style.display = 'block';
        mainPage.style.display = 'none';
        body.style.display = 'flex';
        body.style.alignItems = 'center';
        body.style.justifyContent = 'center';
        body.style.height = '100vh';
      }
    } catch (error) {
      console.error('Error checking session:', error);
    }
  }

  function showMessagingUI() {
    const messagingContainer = document.getElementById('messaging-container');
    if (messagingContainer) {
        messagingContainer.style.display = 'none'; // Initially hidden until user clicks on a contact
    }
    loadUsers(); // Load users list
    connectWebSocket(); // Initialize WebSocket connection
}

  async function logout() {
    try {
      const response = await fetch('/logout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });

      if (response.ok) {
        window.location.reload();
      } else {
        console.error('Logout failed');
      }
    } catch (error) {
      console.error('Error during logout:', error);
    }
  }

  async function handleLogin() {
    const formData = new FormData(loginForm);
    const urlEncodedData = new URLSearchParams(formData);

    try {
      const response = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: urlEncodedData.toString(),
      });

      const result = await response.json();
      if (result.success) {
        window.location.reload();
      } else {
        alert("Login failed: " + result.message);
      }
    } catch (err) {
      alert("Something went wrong.");
    }
  }

  async function createPost() {
    const checkbox = document.getElementById('Create');
    checkbox.checked = false;

    const titleEl = document.getElementById('title');
    const descEl = document.getElementById('description');
    const topicCheckbox = document.querySelectorAll('input[name="topic"]:checked');

    const title = titleEl.value;
    const description = descEl.value;
    const topics = Array.from(topicCheckbox).map(cb => cb.value);

    try {
      const response = await fetch('/posts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, topics })
      });

      if (!response.ok) {
        throw new Error('Server returned status ' + response.status);
      }

      document.getElementById('errorMsg').style.display = 'none';
      await loadPosts();
    } catch (err) {
      console.error('Error creating post:', err);
      const errEl = document.getElementById('errorMsg');
      errEl.textContent = 'An error occurred. Please try again.';
      errEl.style.display = 'block';
    }
  }

  // Post and comment related functions
  async function loadPosts() {
    try {
      const response = await fetch('/api/posts');
      if (!response.ok) throw new Error('Failed to fetch posts');

      const posts = await response.json();
      const postContainer = document.getElementById('Post');
      postContainer.innerHTML = '';

      if (posts.length === 0) {
        postContainer.innerHTML = '<p class="no-posts">No posts yet!</p>';
        return;
      }

      posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.innerHTML = renderPostWithComments(post);
        postContainer.appendChild(postElement);

        loadComments(post.id);
        setupCommentForm(post.id);
      });

      setupReactionListeners();
      setupToggleCommentsListeners();
    } catch (err) {
      console.error('Error loading posts:', err);
      const postContainer = document.getElementById('Post');
      postContainer.innerHTML = '<p class="error">Error loading posts. Please try again later.</p>';
    }
  }

  function renderPostWithComments(post) {
    return `
      <article class="post-card" id="post-${post.id}">
        <div class="post-header">
          <h2>${post.title}</h2>
          <span class="post-meta">Posted by ${post.author}</span>
        </div>
        <div class="post-content">
          <p>${post.content}</p>
        </div>
        <div class="post-footer">
          <div class="topics">
            ${post.topics.map(topic => `<span class="topic-tag">${topic}</span>`).join('')}
          </div>
          <div class="reaction-buttons">
            <button class="like-button" data-post-id="${post.id}" data-type="like">
              <i class="fa fa-thumbs-up"></i> Like (<span class="like-count">${post.likes || 0}</span>)
            </button>
            <button class="dislike-button" data-post-id="${post.id}" data-type="dislike">
              <i class="fa fa-thumbs-down"></i> Dislike (<span class="dislike-count">${post.dislikes || 0}</span>)
            </button>
          </div>
          <span class="post-date">${new Date(post.created_at).toLocaleString()}</span>
        </div>
        <button class="toggle-comments-btn" data-post-id="${post.id}">Show Comments</button>
        <div class="post-comments" id="post-comments-${post.id}" style="display:none;">
          <h3>Comments</h3>
          <form id="comment-form-${post.id}" class="comment-form">
            <textarea placeholder="Add a comment..." required></textarea>
            <button type="submit">Post Comment</button>
          </form>
          <div id="comments-${post.id}" class="comments-container"></div>
        </div>
      </article>
    `;
  }

  function setupReactionListeners() {
    document.querySelectorAll('.like-button, .dislike-button').forEach(button => {
      button.addEventListener('click', async (e) => {
        const btn = e.target.closest('button');
        const postId = btn.dataset.postId;
        const reactionType = btn.dataset.type;
        await handleReaction(postId, reactionType, btn);
      });
    });
  }

  async function handleReaction(postId, reactionType, button) {
    try {
      const response = await fetch('/like', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          reaction_type: reactionType,
          post_id: parseInt(postId)
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        if (response.status === 409) {

        } else {
          throw new Error(data.error || `Failed to ${reactionType} post`);
        }
        return;
      }

      const likeCountSpan = button.parentElement.querySelector('.like-count');
      const dislikeCountSpan = button.parentElement.querySelector('.dislike-count');
      likeCountSpan.textContent = data.like_count;
      dislikeCountSpan.textContent = data.dislike_count;
    } catch (err) {
      console.error(`Error ${reactionType}ing post:`, err);
      alert(`Error ${reactionType}ing post. Please try again.`);
    }
  }

  async function loadComments(postId) {
    const commentsContainer = document.getElementById(`comments-${postId}`);

    try {
      commentsContainer.innerHTML = '<div class="loading">Loading comments...</div>';

      const response = await fetch(`/api/comments?post_id=${postId}`);
      if (!response.ok) throw new Error('Failed to fetch comments');

      const comments = await response.json();
      commentsContainer.innerHTML = '';

      if (comments.length === 0) {
        commentsContainer.innerHTML = `
          <div class="no-comments">
            <i class="fa fa-comment"></i>
            <p>No comments yet. Be the first to comment!</p>
          </div>
        `;
        return;
      }

      comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = 'comment-card';
        commentElement.innerHTML = `
          <div class="comment-header">
            <span class="comment-author">${comment.author}</span>
            <span class="comment-date">${new Date(comment.created_at).toLocaleString()}</span>
          </div>
          <div class="comment-content">
            <p>${comment.content}</p>
          </div>
          <div class="comment-footer">
            <div class="reaction-buttons">
              <button class="comment-like-button" data-comment-id="${comment.id}" data-type="like">
                <i class="fa fa-thumbs-up"></i> Like (<span class="comment-like-count">${comment.likes || 0}</span>)
              </button>
              <button class="comment-dislike-button" data-comment-id="${comment.id}" data-type="dislike">
                <i class="fa fa-thumbs-down"></i> Dislike (<span class="comment-dislike-count">${comment.dislikes || 0}</span>)
              </button>
            </div>
          </div>
        `;
        commentsContainer.appendChild(commentElement);
      });

      setupCommentReactionListeners();
    } catch (err) {
      console.error('Error loading comments:', err);
      commentsContainer.innerHTML = `
        <div class="error-message">
          <i class="fa fa-exclamation-triangle"></i>
          <p>${err.message || 'Error loading comments. Please try again later.'}</p>
          <button class="retry-button" onclick="loadComments(${postId})">Retry</button>
        </div>
      `;
    }
  }

  function setupCommentReactionListeners() {
    document.querySelectorAll('.comment-like-button, .comment-dislike-button').forEach(button => {
      button.addEventListener('click', async (e) => {
        const btn = e.target.closest('button');
        const commentId = btn.dataset.commentId;
        const reactionType = btn.dataset.type;
        await handleCommentReaction(commentId, reactionType, btn);
      });
    });
  }

  async function handleCommentReaction(commentId, reactionType, button) {
    try {
      const response = await fetch('/comment-reaction', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          comment_id: parseInt(commentId),
          reaction_type: reactionType
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        if (response.status === 409) {
          alert("Reaction updated or removed");
        } else {
          throw new Error(data.error || `Failed to ${reactionType} comment`);
        }
        return;
      }

      const likeCountSpan = button.parentElement.querySelector('.comment-like-count');
      const dislikeCountSpan = button.parentElement.querySelector('.comment-dislike-count');
      likeCountSpan.textContent = data.like_count;
      dislikeCountSpan.textContent = data.dislike_count;
    } catch (err) {
      console.error(`Error ${reactionType}ing comment:`, err);
      alert(`Error ${reactionType}ing comment. Please try again.`);
    }
  }

  async function setupCommentForm(postId) {
    const form = document.getElementById(`comment-form-${postId}`);
    const textarea = form.querySelector('textarea');
    const submitButton = form.querySelector('button');
    const errorDiv = document.createElement('div');
    errorDiv.className = 'form-error';
    form.appendChild(errorDiv);

    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const content = textarea.value.trim();

      if (!content) {
        errorDiv.textContent = 'Please enter a comment';
        errorDiv.style.display = 'block';
        return;
      }

      try {
        submitButton.disabled = true;
        submitButton.innerHTML = '<i class="fa fa-spinner fa-spin"></i> Posting...';
        errorDiv.style.display = 'none';

        const response = await fetch('/comments', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            post_id: postId,
            content: content
          }),
        });

        if (!response.ok) throw new Error('Failed to create comment');

        textarea.value = '';
        await loadComments(postId);
      } catch (err) {
        console.error('Error creating comment:', err);
        errorDiv.textContent = err.message || 'Error creating comment. Please try again.';
        errorDiv.style.display = 'block';
      } finally {
        submitButton.disabled = false;
        submitButton.textContent = 'Post Comment';
      }
    });
  }

  function setupToggleCommentsListeners() {
    document.querySelectorAll('.toggle-comments-btn').forEach(button => {
      button.addEventListener('click', function () {
        const postId = button.dataset.postId;
        const commentsDiv = document.getElementById(`post-comments-${postId}`);
        if (commentsDiv.style.display === 'none') {
          commentsDiv.style.display = 'block';
          button.textContent = 'Hide Comments';
        } else {
          commentsDiv.style.display = 'none';
          button.textContent = 'Show Comments';
        }
      });
    });
  }

  // --- Messaging UI and Logic ---

  // Add a simple messaging UI to the main page
  // Add messaging UI to the main page
  const messagingSection = document.createElement('section');
  messagingSection.id = 'messagingSection';
  messagingSection.innerHTML = `
    <div id="messaging-container" style="display:none; margin-top: 30px;">
      <h3 id="user" >Send Message to</h3>
      <form id="messageForm">
       <div id="messagesBorder" class="messages-border">
        <div id="messagesList" class="messages-list"></div>
      </div>
        <input type="text" id="recipient" placeholder="Recipient username" hidden />
        <input type="text" id="messageInput" placeholder="Type your message..." required />
        <button type="submit">Send</button>
      </form>
     
    </div>
  `;
  mainPage.appendChild(messagingSection);

  // Show messaging UI when logged in


  // Fetch and display messages between current user and another user
  async function loadMessages(recipient) {
    try {
      const messagesList = document.getElementById('messagesList');
      messagesList.innerHTML = '<div class="loading">Loading messages...</div>';
      
      const response = await fetch(`/api/messages?recipient=${encodeURIComponent(recipient)}`);
      if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
      }
      const messages = await response.json();
      
      messagesList.innerHTML = '';
      if (messages.length === 0) {
          messagesList.innerHTML = '<p class="no-messages">No messages yet.</p>';
          return;
      }

      messages.forEach(msg => {
          const msgDiv = document.createElement('div');
          msgDiv.className = `message-item ${msg.sender === currentUser ? 'sent' : 'received'}`;
          msgDiv.innerHTML = `
              <div class="message-content">
                  <span class="message-sender">${msg.sender}</span>
                  <p>${msg.content}</p>
                  <span class="message-time">${new Date(msg.timestamp).toLocaleString()}</span>
              </div>
          `;
          messagesList.appendChild(msgDiv);
      });
      
      // Scroll to bottom of messages
      messagesList.scrollTop = messagesList.scrollHeight;
    } catch (err) {
      console.error('Error loading messages:', err);
      document.getElementById('messagesList').innerHTML = 
          '<p class="error">Error loading messages. Please try again.</p>';
    }
  }

  // Handle sending a message
  const messageForm = document.getElementById('messageForm');
  if (messageForm) {
    messageForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const recipient = document.getElementById('recipient').value.trim();
      const content = document.getElementById('messageInput').value.trim();
      if (!recipient || !content) return;
      try {
        const response = await fetch('/messages', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ recipient, content })
        });
        if (!response.ok) throw new Error('Failed to send message');
        document.getElementById('messageInput').value = '';
        await loadMessages(recipient);
      } catch (err) {
        alert('Error sending message.');
      }
    });
  }

  // Optionally, you can add a button to open the messaging UI and load messages for a specific user
  // For now, show the messaging UI when logged in
  checkSession = async function () {
    try {
      const response = await fetch('/check-session');
      const data = await response.json();
      if (data.loggedIn) {
        mainPage.style.display = 'block';
        usernameDisplay.textContent = `Welcome, ${data.username}!`;
        container.style.display = 'none';
        loadPosts();
        showMessagingUI();
      } else {
        container.style.display = 'block';
        mainPage.style.display = 'none';
        body.style.display = 'flex';
        body.style.alignItems = 'center';
        body.style.justifyContent = 'center';
        body.style.height = '100vh';
      }
    } catch (error) {
      console.error('Error checking session:', error);
    }
  }

  // Add user sidebar to the main page
  const userSidebar = document.createElement('div');
  userSidebar.id = 'userSidebar';

  userSidebar.innerHTML = '<h3 style="text-align:center;color:#00ff9d;margin-bottom:10px;">Users</h3><div id="userList"></div>';
  document.body.appendChild(userSidebar);

  // Fetch and display all users in the sidebar
  async function loadUsers() {
    try {
      const response = await fetch('/api/users');
      if (!response.ok) throw new Error('Failed to fetch users');
      const users = await response.json();
      const userList = document.getElementById('userList');
      userList.innerHTML = '';
      users.forEach(user => {
        const userDiv = document.createElement('div');
        userDiv.className = 'user-list-item';
        userDiv.style = 'padding:6px 0;cursor:pointer;border-bottom:1px solid #eee;';
        userDiv.textContent = user;
        userDiv.addEventListener('click', () => {
          document.getElementById('recipient').value = user;
          document.getElementById('user').textContent = `Send Message to ${user}`;
          loadMessages(user);
          if(document.getElementById('messaging-container').style.display == 'block'){
            document.getElementById('messaging-container').style.display = 'none';
          }else{
            document.getElementById('messaging-container').style.display = 'block';
          }
        });
        userList.appendChild(userDiv);
      });
    } catch (err) {
      document.getElementById('userList').innerHTML = '<p>Error loading users.</p>';
    }
  }

  // Call loadUsers on page load
  loadUsers();

  let ws = null;
  let currentChatUser = null;

  function connectWebSocket() {
      ws = new WebSocket(`ws://${window.location.host}/ws`);
      
      ws.onopen = () => {
          console.log('WebSocket Connected');
      };

      ws.onmessage = (event) => {
          const message = JSON.parse(event.data);
          
          switch(message.type) {
              case 'userList':
                  updateOnlineUsers(message.content);
                  break;
              case 'privateMessage':
                  handleNewMessage(message);
                  break;
          }
      };

      ws.onclose = () => {
          console.log('WebSocket Disconnected');
          setTimeout(connectWebSocket, 1000);
      };
  }

  function updateOnlineUsers(users) {
      const userList = document.getElementById('userList');
      if (!userList) return;
      
      userList.innerHTML = users.map(user => `
          <div class="user-list-item ${user.online ? 'online' : 'offline'}"
               onclick="openChat('${user.username}')">
              ${user.username}
              ${user.online ? '🟢' : '⚪'}
          </div>
      `).join('');
  }

  function handleNewMessage(message) {
    const messagesList = document.getElementById('messagesList');
    if (!messagesList) return;

    const msgDiv = document.createElement('div');
    msgDiv.className = 'message-item ' + (message.from === currentUser ? 'sent' : 'received');
    msgDiv.innerHTML = `
        <div class="message-content">
            <p>${message.content}</p>
            <span class="message-time">${new Date(message.timestamp).toLocaleString()}</span>
        </div>
    `;
    messagesList.appendChild(msgDiv);
    messagesList.scrollTop = messagesList.scrollHeight;
}

// Add message pagination
let messageOffset = 0;
const messageLimit = 10;

async function loadMoreMessages(recipient) {
    try {
        const response = await fetch(`/api/messages?recipient=${encodeURIComponent(recipient)}&offset=${messageOffset}&limit=${messageLimit}`);
        if (!response.ok) throw new Error('Failed to fetch messages');
        const messages = await response.json();
        const messagesList = document.getElementById('messagesList');
        messages.forEach(msg => {
            const msgDiv = document.createElement('div');
            msgDiv.className = 'message-item';
            msgDiv.innerHTML = `<b>${msg.sender}:</b> ${msg.content} <span style="font-size:0.8em;color:#888;">${new Date(msg.timestamp).toLocaleString()}</span>`;
            messagesList.appendChild(msgDiv);
        });
        messageOffset += messages.length;
    } catch (err) {
        console.error('Error loading messages:', err);
    }
}

  function sendMessage(content) {
      if (!ws || !currentChatUser) return;
      
      ws.send(JSON.stringify({
          type: 'privateMessage',
          to: currentChatUser,
          content: content
      }));
  }

  // Update message form handler
  messageForm?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const content = document.getElementById('messageInput').value.trim();
      if (!content) return;
      
      sendMessage(content);
      document.getElementById('messageInput').value = '';
  });

  // Initialize WebSocket connection when logged in
  if (mainPage.style.display === 'block') {
      connectWebSocket();
  }
});