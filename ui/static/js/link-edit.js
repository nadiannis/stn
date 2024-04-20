const linkEditForm = document.getElementById('link-edit-form');
const errors = linkEditForm.querySelectorAll('.error');
const errorContainers = {};
errors.forEach((errorContainer) => {
  const key = errorContainer.parentElement.classList[0];
  errorContainer.style.display = 'none';
  errorContainers[key] = errorContainer;
});

linkEditForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const id = location.pathname.split('/')[2];
  for (const [key, errorContainer] of Object.entries(errorContainers)) {
    e.target[key].classList.remove('input-error');
    errorContainer.textContent = '';
    errorContainer.style.display = 'none';
  }

  try {
    const response = await fetch(`${location.origin}/links/${id}/edit`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams({
        url: e.target['url'].value,
        'back-half': e.target['back-half'].value,
      }),
    });

    if (!response.ok) {
      const err = new Error(`${response.status} ${response.statusText}`);
      err.status = response.status;
      err.response = response;
      throw err;
    }

    location.href = response.url;
  } catch (error) {
    if (error.response?.headers.get('Content-Type') === 'application/json') {
      const result = await error.response.json();
      console.log(result);

      for (const [key, val] of Object.entries(result.fieldErrors)) {
        e.target[key].classList.add('input-error');
        errorContainers[key].textContent = val;
        errorContainers[key].style.display = 'inline';
      }
    }
  }
});
