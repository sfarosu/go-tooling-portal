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

function validateUserPass() {
  var x = document.forms["htpasswd"]["username"].value;
  var y = document.forms["htpasswd"]["password"].value;
  if (x == "" || y == "") {
    alert("Username or password field cannot be empty !");
    return false;
  }
  if (y.length < 8 || y.length > 64) {
    alert("Password length must contain between 8 and 64 characters !");
    return false;
  }
}

function validatePassgenSelection(){
  var length = document.forms["passgen"]["number"].value;
  if (length < 4 || length > 64) {
    alert("Password length must contain between 4 and 64 characters !");
    return false;
  }
  var uppercase = document.getElementById("uppercase");
  var lowercase = document.getElementById("lowercase");
  var numbers = document.getElementById("numbers");
  var symbols = document.getElementById("symbols");
  if (uppercase.checked == false && lowercase.checked == false && numbers.checked == false && symbols.checked == false ) {
    alert("You must select at least a category of characters !");
    return false;
  }
}

function validateSshEmailPass() {
  var x = document.forms["sshkeygen"]["email"].value;
  var y = document.forms["sshkeygen"]["pass"].value;
  if (x.length === 0) {
    alert("You must enter an e-mail address !");
    return false;
  }
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

  if (x.length === 0) {
    alert("Textbox is empty !");
    return false;
  }

  try {
    JSON.parse(x);
  } catch (e) {
    alert("Json is not valid !");
    return false;
  }
  return true;
}

function validateConvert() {
  var x = document.forms["formatconvert"]["text"].value;

  if (x.length === 0) {
    alert("Textbox is empty !");
    return false;
  }
  return true;
}

function copyToClipboard(elemId) {
  var copyText = document.getElementById(elemId);
  copyText.select();
  navigator.clipboard.writeText(copyText.value);
}

// function validateFileSelection(file) {
//   var selectedfile = document.getElementById(file).value

//   if (selectedfile == "") {
//     alert("You must select a file !");
//     return false;
//   }
// }
