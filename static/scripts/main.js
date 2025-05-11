document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const loginForm = document.getElementById('loginForm');
  const mainPage = document.getElementById('mainPage');
  const logoutBtn = document.getElementById('logout');
  const usernameDisplay = document.getElementById('username'); // Where we display the username

  // Check session on load
  fetch('/check-session')
    .then(res => res.json())
    .then(data => {
      if (data.loggedIn) {
        container.style.display = 'none';
        mainPage.style.display = 'block';
        document.body.style.display = 'block';
        usernameDisplay.textContent = `Welcome, ${data.username}!`; 
        
      } else {
        container.style.display = 'block';
        mainPage.style.display = 'none';
        document.body.style.display = 'flex';
      }
    });

  // Logout functionality
  logoutBtn.addEventListener('click', async () => {
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
  });

  // Toggle between Sign Up and Sign In forms
  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    loginForm.style.display = 'none';
  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    loginForm.style.display = 'block';
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
