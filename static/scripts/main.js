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
  document.getElementById('form')?.addEventListener('submit', async function(e) {
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
        usernameDisplay.textContent = `Welcome, ${data.username}!`;
        container.style.display = 'none';
        loadPosts();
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
        <div class="post-comments">
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
          alert("Reaction updated or removed");
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
});