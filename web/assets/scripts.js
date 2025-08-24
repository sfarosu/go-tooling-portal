/**
 * Initializes the search functionality for filtering and highlighting cards.
 *
 * - Finds the search input field with id "searchInput".
 * - Listens for input events on the search box.
 * - On each input event:
 *   - Converts the input value to lowercase for case-insensitive matching.
 *   - Iterates over all elements with the class "card".
 *   - For each card, checks if the card's title text includes the search filter.
 *     - If matched and filter is not empty:
 *       - Shows the card by setting its container's display style to default.
 *       - Highlights all occurrences of the matched text inside the card title using <mark> tags.
 *     - If not matched:
 *       - Hides the card by setting its container's display style to "none".
 *       - Removes any existing highlights from the card title.
 *   - If the input is empty:
 *       - Shows all cards and removes all highlights to reset the view.
 */
function initSearch() {
  const input = document.getElementById('searchInput');
  if (!input) return;

  input.addEventListener('input', () => {
    const filter = input.value.toLowerCase();
    const cards = document.querySelectorAll('.card');

    cards.forEach(card => {
      const titleEl = card.querySelector('.card-title');
      const titleText = titleEl.textContent;
      const titleLower = titleText.toLowerCase();

      if (titleLower.includes(filter) && filter !== '') {
        // Show card column
        card.parentElement.style.display = '';

        // Highlight matched part(s)
        // Escape special regex characters in filter for safe usage
        const escapedFilter = filter.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

        // Create regex, global and case-insensitive
        const regex = new RegExp(`(${escapedFilter})`, 'gi');

        // Replace matched substring(s) with <mark>
        const highlighted = titleText.replace(regex, '<mark>$1</mark>');
        titleEl.innerHTML = highlighted;

      } else {
        // Hide card column
        card.parentElement.style.display = 'none';

        // Remove any existing highlights when hiding
        titleEl.innerHTML = titleText;
      }

      // If input is empty, reset all highlights and show all cards
      if (filter === '') {
        card.parentElement.style.display = '';
        titleEl.innerHTML = titleText;
      }
    });
  });
}

/*
 * Sets up an event listener to detect when HTMX swaps in the navbar content.
 * When the navbar is loaded (inside #navbar-container), it calls initSearch()
 * to attach the search input event handlers.
 */
function setupNavbarSearchInit() {
  document.body.addEventListener('htmx:afterSwap', (event) => {
    // Only trigger when the navbar-container was swapped
    if (event.detail.target.id === 'navbar-container') {
      if (typeof initSearch === 'function') {
        initSearch();
      }
    }
  });
}

// When the initial page DOM content is loaded, setup the HTMX listener for navbar loading
document.addEventListener('DOMContentLoaded', () => {
  setupNavbarSearchInit();
});

// Reusable password visibility toggle
function togglePasswordVisibility(passwordInputId, iconId) {
    const passwordInput = document.getElementById(passwordInputId);
    const icon = document.getElementById(iconId);

    if (!passwordInput) return;

    const isVisible = passwordInput.type === 'text';
    passwordInput.type = isVisible ? 'password' : 'text';

    if (icon) {
        icon.className = isVisible ? 'bi bi-eye-slash' : 'bi bi-eye';
    }
}

// Function to copy text from an input element to the clipboard
async function copyToClipboard(elemId) {
  try {
    const element = document.getElementById(elemId);
    if (!element) {
      console.warn(`Element with id="${elemId}" not found.`);
      return;
    }

    // Get text depending on element type
    const tag = element.tagName.toUpperCase();
    const text = (tag === 'INPUT' || tag === 'TEXTAREA')
      ? element.value
      : element.innerText || element.textContent;

    await navigator.clipboard.writeText(text);

    // Visual feedback: highlight the element
    element.classList.add('copied-highlight');
    setTimeout(() => {
      element.classList.remove('copied-highlight');
    }, 700);

  } catch (err) {
    console.error('Failed to copy text: ', err);
  }
}

// Function to initialize all Bootstrap tooltips on the page
function initBootstrapTooltips() {
  const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
  tooltipTriggerList.forEach(function (tooltipTriggerEl) {
    new bootstrap.Tooltip(tooltipTriggerEl);
  });
}
// Call it when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  initBootstrapTooltips();
  setupNavbarSearchInit();
});