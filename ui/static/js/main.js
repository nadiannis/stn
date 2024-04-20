lucide.createIcons();

const backButton = document.querySelector('.btn-back');

if (backButton) {
  backButton.addEventListener('click', () => {
    history.back();
  });
}

const deleteLink = async (id) => {
  try {
    const response = await fetch(`${location.origin}/links/${id}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      const err = new Error(`${response.status} ${response.statusText}`);
      err.status = response.status;
      err.response = response;
      throw err;
    }

    location.href = response.url;
  } catch (error) {
    console.log(error);
  }
};
