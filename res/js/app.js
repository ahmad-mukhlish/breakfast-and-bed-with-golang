function notify(msg, msgType) {
    notie.alert({
        type: msgType,
        text: msg,
        position: "bottom"
    })
}

function sweetAlert(title, html, icon, confirmButtonText) {
    Swal.fire({
        title: title,
        html: html,
        icon: icon,
        confirmButtonText: confirmButtonText,
    })
}

function prompt() {

    let toast = function (param) {
        const Toast = Swal.mixin({
            toast: true,
            position: param.pos ?? "top-end",
            showConfirmButton: false,
            timer: 3000,
            icon: param.icon ?? "success",
            title: param.msg ?? "",
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });
        Toast.fire();
    }

    let success = function (title, msg) {
        Swal.fire({
            title: title ?? "Success",
            text: msg ?? "You're awesome",
            icon: "success"
        });
    }

    let error = function (title, msg) {
        Swal.fire({
            title: title ?? "Success",
            text: msg ?? "You're awesome",
            icon: "error"
        });
    }

    let inputDates = async function (callback) {


        let htmlDate = `

<form name="pickerInside" action="" class="mt-3 needs-validation px-3" novalidate id="pickerInsidePopup">

<div class="row" id="range-picker-inside">
<div class="col">

<div class="mb-3">
    <input disabled type="text" class="form-control" id="start-inside" name="start" aria-describedby="startDate" required placeHolder="Arrival" autocomplete="off">
    <div class="invalid-feedback">
        Please input a date
    </div>
</div>

</div>

<div class="col">

<div class="mb-3">
    <input disabled type="text" class="form-control" id="end-inside" name="end" aria-describedby="endDate" required placeHolder="Departure" autocomplete="off">
    <div class="invalid-feedback">
        Please input a date
    </div>
</div>

</div>
</div>

<button type="submit" class="btn btn-primary mt-3">Check for availability</button>
</form>
`;

        await Swal.fire({
            title: "Search for Availability",
            html: htmlDate,
            confirmButtonText: "Check for Availabilty",
            showConfirmButton: false,
            willOpen: () => {

                const elem = document.getElementById("range-picker-inside");
                const rangepicker = new DateRangePicker(elem, {
                    orientation: "top auto",
                    buttonClass: "btn",
                    minDate: new Date(),
                    format : "dd-mm-yyyy"
                });


                const forms = document.querySelectorAll('.needs-validation');
                Array.from(forms).forEach(form => {

                    form.addEventListener('submit', function (event) {
                        if (!form.checkValidity()) {
                            event.preventDefault();
                            event.stopPropagation();

                        } else {
                            event.preventDefault();
                            event.stopPropagation();
                            Swal.clickConfirm();
                        }

                        form.classList.add('was-validated');


                    });


                });
            },

            didOpen: () => {
                const start = document.getElementById("start-inside");
                start.removeAttribute("disabled");

                const end = document.getElementById("end-inside");
                end.removeAttribute("disabled");

                // document.pickerInside.submit();

            },

            preConfirm: () => {
                return [
                    document.getElementById("start-inside").value,
                    document.getElementById("end-inside").value,
                ];

            },


        }).then(result => {

            if (callback === undefined) {
                return;
            }

            if (result === null) {
                callback("");
                return;
            }

            if (result.dismiss === Swal.DismissReason.cancel) {
                callback("");
                return;
            }

            callback(result);

        });

    }

    return {
        toast: toast,
        success: success,
        error: error,
        inputDates: inputDates,
    }

}


function showAvailabilityForRoomsById(roomId) {
    let attention = prompt();

    let button = document.getElementById("check");

    button.addEventListener("click", function () {
        attention.inputDates(_ => {

            let form = document.getElementById("pickerInsidePopup");
            let formData = new FormData(form)
            let tokenElem = document.getElementById("token")
            formData.append("csrf_token", tokenElem.value)
            formData.append("room_id", roomId)

            fetch('/check-availability/json', {
                method: "POST",
                body: formData,
            })
                .then(response => response.json())
                .then(jsonData => Swal.fire({
                    title: jsonData.message,
                }).then((result) => {
                    window.open("/check/room/"+roomId, "_self")
                }));
        });
    });
}
