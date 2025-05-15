document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');


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
      } else {
        container.style.display = 'block';
        mainPage.style.display = 'none';
        body.style.display = 'flex';
        body.style.alignItems = 'center';
        body.style.justifyContent = 'center';
        body.style.height='100vh'

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
});
