function revealPassword(elemId) {
	var x = document.getElementById(elemId);
	if (x.type === "password") {
	  x.type = "text";
	} else {
	  x.type = "password";
	}
}

function deactivateField(elemId) {
  if(document.getElementById('no').selected == true){
    document.getElementById(elemId).hidden = true;
  }
  if(document.getElementById('yes').selected == true){
    document.getElementById(elemId).hidden = false;
  }
}

function validatePassgenSelection(){
  var uppercase = document.getElementById("uppercase");
  var lowercase = document.getElementById("lowercase");
  var numbers = document.getElementById("numbers");
  var symbols = document.getElementById("symbols");
  if (uppercase.checked == false && lowercase.checked == false && numbers.checked == false && symbols.checked == false ) {
    alert("You must select at least a category of characters !");
    return false;
  }
}

function handleSshPasswordField() {
  var y = document.forms["sshkeygen"]["pass"].value;

  if (document.getElementById('no').selected == false) {
    if (y.length < 8 || y.length > 64)  {
      alert("Password length must contain between 8 and 64 characters !");
      return false;
    }
  } else {
    document.getElementById("pass").value = ""; //we make sure that if the user selects yes, inserts few chars and then selects no, the password field is empty
  }
}

function validateJson() {
  var x = document.forms["jsonprettify"]["text"].value;

  try {
    JSON.parse(x);
  } catch (e) {
    alert("Json is not valid !");
    return false;
  }
  return true;
}

function copyToClipboard(elemId) {
  var copyText = document.getElementById(elemId);
  copyText.select();
  navigator.clipboard.writeText(copyText.value);
}

function getClientTimeZone() {
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  document.getElementById("browserTimeZoneFromEpochForm").value = timezone;
  document.getElementById("browserTimeZoneFromHumanForm").value = timezone;
}

function changeEpochTimeFormat(inputField, selectorField) {
  document.getElementById(inputField).value = document.getElementById(selectorField).value
}

function disableEpochTimeFormatSelect() {
  var inputValue = document.getElementById("epochtime").value;
  if (inputValue != '') {
    document.getElementById("epochTimeFormat").setAttribute("disabled", "disabled");
  }
}

function getSelectedTimeConversion(sourceForm) {
  if (sourceForm == 'epochform') {
    document.getElementById("epochToHuman").value = "true";
  } else if (sourceForm == 'humanform') {
    document.getElementById("humanToEpoch").value = "true";
  }
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

function downloadFile(content, filename) {
  if (!content) {
      console.error('Content is empty or not provided.');
      return;
  }

  var blob = new Blob([content], { type: 'application/octet-stream' });
  var url = URL.createObjectURL(blob);
  var a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
