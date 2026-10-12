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
    const target = event.detail.target;
    if (target.querySelector('.dropdown-menu')) {
      initActiveToolLink(target);
    }

    // Only trigger when the navbar-container was swapped
    if (target.id === 'navbar-container') {
      if (typeof initSearch === 'function') {
        initSearch();
      }
    }
  });
}

function initActiveToolLink(navbar = document) {
  const currentPath = window.location.pathname;
  navbar.querySelectorAll('.dropdown-menu a.dropdown-item[href]').forEach(link => {
    const isCurrentPage = new URL(link.href, window.location.href).pathname === currentPath;
    link.classList.toggle('active', isCurrentPage);
    if (isCurrentPage) {
      link.setAttribute('aria-current', 'page');
    } else {
      link.removeAttribute('aria-current');
    }
  });
}

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

function initTimeConverterTimezone() {
  const timezoneSelect = document.getElementById('timezone-select');
  const timezoneLabel = document.getElementById('timezone-label');
  if (!timezoneSelect || !timezoneLabel) return;

  const fallbackTimezones = [
    'UTC',
    'America/Los_Angeles',
    'America/Denver',
    'America/Chicago',
    'America/New_York',
    'Europe/London',
    'Europe/Paris',
    'Asia/Kolkata',
    'Asia/Tokyo',
    'Australia/Sydney'
  ];
  let supportedTimezones = fallbackTimezones;
  if (typeof Intl.supportedValuesOf === 'function') {
    try {
      const availableTimezones = Intl.supportedValuesOf('timeZone');
      if (availableTimezones.length > 0) supportedTimezones = availableTimezones;
    } catch (error) {
      console.warn('Could not load the browser timezone list; using common timezones.', error);
    }
  }

  let browserTimezone = 'UTC';
  try {
    browserTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
  } catch (error) {
    console.warn('Could not detect the browser timezone; using UTC.', error);
  }

  const timezones = new Set(['UTC', ...supportedTimezones, browserTimezone]);
  [...timezones].sort((first, second) => first.localeCompare(second)).forEach(timezone => {
    const option = document.createElement('option');
    option.value = timezone;
    option.textContent = timezone;
    timezoneSelect.appendChild(option);
  });

  timezoneLabel.textContent = browserTimezone;
  timezoneSelect.value = browserTimezone;

  const epochInput = document.getElementById('epochtime');
  const epochFormatSelect = document.getElementById('epochTimeFormat');
  const epochTimeHelp = document.getElementById('epoch-time-help');
  const epochFormats = {
    seconds: {
      pattern: '[0-9]{1,10}',
      title: 'Enter up to 10 digits for a seconds timestamp',
      help: 'Up to 10 digits for seconds'
    },
    millis: {
      pattern: '[0-9]{1,13}',
      title: 'Enter up to 13 digits for a milliseconds timestamp',
      help: 'Up to 13 digits for milliseconds'
    },
    micros: {
      pattern: '[0-9]{1,16}',
      title: 'Enter up to 16 digits for a microseconds timestamp',
      help: 'Up to 16 digits for microseconds'
    }
  };
  const updateEpochTimeFormat = (now = new Date()) => {
    if (!epochInput || !epochFormatSelect) return;

    const format = epochFormats[epochFormatSelect.value];
    if (!format) return;

    epochInput.pattern = format.pattern;
    epochInput.title = format.title;
    const precisionInput = document.getElementById('epoch-precision');
    if (precisionInput) precisionInput.value = epochFormatSelect.value;
    if (epochTimeHelp) epochTimeHelp.textContent = format.help;

    const nowInMilliseconds = now.getTime();
    const timestamp = {
      seconds: Math.floor(nowInMilliseconds / 1000),
      millis: nowInMilliseconds,
      micros: nowInMilliseconds * 1000
    }[epochFormatSelect.value];
    epochInput.value = timestamp.toString();
  };

  const humanTimeFieldIds = ['year', 'month', 'day', 'hour', 'minute', 'second', 'millisecond', 'microsecond'];
  const humanTimeFields = humanTimeFieldIds.map(id => document.getElementById(id)).filter(Boolean);
  let humanTimeEdited = false;
  humanTimeFields.forEach(field => {
    field.addEventListener('input', () => {
      humanTimeEdited = true;
    });
  });

  const updateHumanTimeToNow = (timezone, now = new Date()) => {
    const parts = new Intl.DateTimeFormat('en-US', {
      timeZone: timezone,
      year: 'numeric',
      month: 'numeric',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
      hourCycle: 'h23'
    }).formatToParts(now);
    const values = Object.fromEntries(parts.map(part => [part.type, part.value]));
    const timeValues = {
      year: values.year,
      month: values.month,
      day: values.day,
      hour: values.hour,
      minute: values.minute,
      second: values.second,
      millisecond: now.getMilliseconds(),
      microsecond: 0
    };

    Object.entries(timeValues).forEach(([id, value]) => {
      const field = document.getElementById(id);
      if (field) field.value = value;
    });
  };

  const updateFormTimezones = () => {
    ['epoch-timezone', 'human-timezone'].forEach(id => {
      const timezoneInput = document.getElementById(id);
      if (timezoneInput) timezoneInput.value = timezoneSelect.value;
    });
    ['epoch-selected-timezone', 'human-selected-timezone'].forEach(id => {
      const timezoneDisplay = document.getElementById(id);
      if (timezoneDisplay) timezoneDisplay.textContent = timezoneSelect.value;
    });
    if (!humanTimeEdited) updateHumanTimeToNow(timezoneSelect.value, referenceTime);
  };

  let referenceTime = new Date();
  const updateCurrentTimestamps = (now = new Date()) => {
    referenceTime = epochFormatSelect?.value === 'seconds'
      ? new Date(Math.floor(now.getTime() / 1000) * 1000)
      : now;
    updateEpochTimeFormat(referenceTime);
    if (!humanTimeEdited) updateHumanTimeToNow(timezoneSelect.value, referenceTime);
  };

  if (epochFormatSelect) {
    epochFormatSelect.addEventListener('change', () => {
      updateCurrentTimestamps(new Date());
      ['epoch-result', 'human-result'].forEach(id => {
        const result = document.getElementById(id);
        if (result) result.replaceChildren();
      });
    });
  }

  timezoneSelect.addEventListener('change', () => {
    updateFormTimezones();
    ['epoch-result', 'human-result'].forEach(id => {
      const result = document.getElementById(id);
      if (result) result.replaceChildren();
    });
  });
  if (epochFormatSelect) updateCurrentTimestamps(referenceTime);
  updateFormTimezones();
}

function initTimeConverterRequestEncoding() {
  if (!document.body) return;

  document.body.addEventListener('htmx:configRequest', event => {
    const trigger = event.detail.elt;
    if (!trigger) return;
    const form = trigger instanceof HTMLFormElement ? trigger : trigger.closest('form');
    if (!form || !['epoch-to-human-form', 'human-to-epoch-form'].includes(form.id)) return;

    const parameters = event.detail.parameters;
    parameters.epochToHumanMode = form.id === 'epoch-to-human-form';
    parameters.humanToEpochMode = form.id === 'human-to-epoch-form';

    ['epochTime', 'year', 'month', 'day', 'hour', 'minute', 'second', 'millisecond', 'microsecond'].forEach(name => {
      parameters[name] = Number(parameters[name] || 0);
    });
  });
}

// Display/hide forms depending on a select value
function toggleForms(selectElementId, optionValues, formGroupIds) {
  var selectElement = document.getElementById(selectElementId);

  // Hide all form groups initially
  formGroupIds.forEach(function(formGroupId) {
      document.getElementById(formGroupId).classList.add('d-none');
  });

  // Determine which form group to display based on the selected option value
  var selectedIndex = optionValues.indexOf(selectElement.value);
  if (selectedIndex !== -1) {
      document.getElementById(formGroupIds[selectedIndex]).classList.remove('d-none');
  }
}

// Call functions when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  initBootstrapTooltips();
  initTimeConverterTimezone();
  initTimeConverterRequestEncoding();
  setupNavbarSearchInit();
});
