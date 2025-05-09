document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const loginForm = document.getElementById('loginForm');
  const mainPage = document.getElementById('mainPage');
  const logoutBtn = document.getElementById('logoutBtn');

  // Check session on load
  fetch('/check-session')
  .then(res => res.json())
  .then(data => {
    if (data.loggedIn) {
      container.style.display = 'none';
      mainPage.style.display = 'block';
    } else {
      container.style.display = 'block';
      mainPage.style.display = 'none';
    }
  });


  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    loginForm.style.display = 'none';
  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    loginForm.style.display = 'block';
  });

  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const formData = new FormData(loginForm);
    const urlEncodedData = new URLSearchParams(formData);

    try {
      const response = await fetch('/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: urlEncodedData.toString(),
      });

      const result = await response.json();
      if (result.success) {
        container.style.display = 'none';
        mainPage.style.display = 'block';
      } else {
        alert("Login failed: " + result.message);
      }
    } catch (err) {
      alert("Something went wrong.");
    }
  });

  logoutBtn?.addEventListener('click', async () => {
    await fetch('/logout', { method: 'POST' });
    container.style.display = 'block';
    mainPage.style.display = 'none';
  });
});
