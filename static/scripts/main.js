document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const section = document.querySelector('section');
  section.innerHTML = `<form id="form"  method="get">
                <div class="container">
                    <h3>Create Post</h3>
                    <div class="div-title">
                        <label for="title">Title :</label>
                        <input type="text" name="title" id="title" required>
                    </div>
                    <div class="div-description">
                        <label for="description">description :</label>
                        <textarea name="description" id="description" rows="4"  required></textarea>
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
            </form>`
;
  // Dynamically add sign-up and sign-in forms
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
              <input type="text" pla  const mainPage = document.getElementById('mainPage');
  const usernameDisplay = document.getElementById('username');
  const logoutBtn = document.getElementById('logout');
  const loginForm = document.getElementById('loginForm');ceholder="First Name" name="first_name" required />
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
  const mainPage = document.getElementById('mainPage');
  const usernameDisplay = document.getElementById('username');
  const logoutBtn = document.getElementById('logout');
  const loginForm = document.getElementById('loginForm');
  // Toggle between Sign Up and Sign In forms
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const body = document.getElementById('body')

  // Toggle between Sign Up and Sign In forms
  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    loginForm.style.display = 'none';
  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    loginForm.style.display = 'block';
  });
  // Check session on load
  fetch('/check-session')
    .then(res => res.json())
    .then(data => {
      if (data.loggedIn) {
        mainPage.style.display = 'block';
        usernameDisplay.textContent = `Welcome, ${data.username}!`;
        container.style.display = 'none';
        // Load posts after successful login
        loadPosts();
      } else {
        container.style.display = 'block';
        mainPage.style.display = 'none';
        body.style.display = 'flex';
        body.style.alignItems = 'center';
        body.style.justifyContent = 'center';
        body.style.height = '100vh'

      }
    });

  // Logout functionality
  logoutBtn.addEventListener('submit', async (e) => {
    e.preventDefault();

    try {
      const response = await fetch('/logout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });

      if (response.ok) {
        window.history.pushState({ action: 'logout' }, 'Logged Out', '/logout');
        window.location.reload();
      } else {
        console.error('Logout failed');
      }
    } catch (error) {
      console.error('Error during logou59pxt:', error);
    }
  });

  // Handle login form submission
  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();

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
  });

  // Add form submission handler inside DOMContentLoaded
  document.getElementById('form').addEventListener('submit', async function(e) {
    e.preventDefault();
    const checkbox = document.getElementById('Create');
    checkbox.checked = false;
   
    const titleEl       = document.getElementById('title');
    const descEl        = document.getElementById('description');
    const topicCheckbox = document.querySelectorAll('input[name="topic"]:checked');


    const title       = titleEl.value;
    const description = descEl.value;
    const topics      = Array.from(topicCheckbox).map(cb => cb.value);

    try {
      const response = await fetch('/posts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, topics })
      });

      if (!response.ok) {
        console.error('Server returned status', response.status);
        return;
      }

      const payload = await response.json();
     console.log('Title:', payload.title);
      console.log('Description:', payload.content);
      console.log('Topics:', payload.interest);
      const post = document.getElementById('Post');
      post.innerHTML = `
        <h2>${payload.title}</h2>
        <p>${payload.content}</p>
        <p>Topics: ${payload.interest}</p>
      `;
      document.getElementById('errorMsg').style.display = 'none';
      await loadPosts(); // Reload posts after successful creation
    } catch (err) {
      console.error('Error sending form:', err);
      const errEl = document.getElementById('errorMsg');
      errEl.textContent = 'An error occurred. Please try again.';
      errEl.style.display = 'block';
    }
  });

  // Move loadPosts function outside the DOMContentLoaded listener
  async function loadPosts() {
    try {
      const response = await fetch('/api/posts');
      if (!response.ok) {
        throw new Error('Failed to fetch posts');
      }
      const posts = await response.json();
      
      const postContainer = document.getElementById('Post');
      postContainer.innerHTML = ''; // Clear existing posts
      
      if (posts.length === 0) {
        postContainer.innerHTML = '<p class="no-posts">No posts yet!</p>';
        return;
      }
      
      posts.forEach(post => {
        const postElement = document.createElement('article');
        postElement.className = 'post-card';
        postElement.innerHTML = `
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
            <span class="post-date">${new Date(post.created_at).toLocaleString()}</span>
          </div>
        `;
        postContainer.appendChild(postElement);
      });
    } catch (err) {
      console.error('Error loading posts:', err);
      const postContainer = document.getElementById('Post');
      postContainer.innerHTML = '<p class="error">Error loading posts. Please try again later.</p>';
    }
  }
});




