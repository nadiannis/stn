const homeLinkCreateForm = document.getElementById('home-link-create-form');
const heroCta = document.querySelector('.hero-cta');

homeLinkCreateForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const resultContainer = heroCta.querySelector('.hero-cta-result');
  const errorContainer = heroCta.querySelector('.error');
  resultContainer.classList.remove('short-link-container');
  resultContainer.innerHTML = '';
  errorContainer.textContent = '';

  try {
    const response = await fetch(
      `${location.origin}/links/create?from=home`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
          url: e.target.url.value,
        }),
      }
    );

    if (!response.ok) {
      const err = new Error(`${response.status} ${response.statusText}`);
      err.status = response.status;
      err.response = response;
      throw err;
    }

    const data = await response.json();
    console.log(data);

    resultContainer.classList.add('short-link-container');
    const a = document.createElement('a');
    a.textContent = `${location.origin}/${data.backHalf}`;
    a.setAttribute('href', `${location.origin}/${data.backHalf}`);
    a.setAttribute('target', '_blank');
    resultContainer.append(a);
    e.target.url.value = '';
  } catch (error) {
    if (error.response.headers.get('Content-Type') === 'application/json') {
      const result = await error.response.json();
      console.log(result);

      errorContainer.textContent = result.fieldErrors.url;
    }
  }
});
