function prompt() {
    let toast = function(content) {
      const {
        msg = "",
        icon = "success",
        position = "top-end"
      } = content;
    
      console.log("creating toast");
      const Toast = Swal.mixin({
        toast: true,
        title: msg,
        icon: icon,
        position: position,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.onmouseenter = Swal.stopTimer;
          toast.onmouseleave = Swal.resumeTimer;
        }
      });

      Toast.fire({});
    }

    let success = function(content) {
      const {
        msg = "",
        title = "",
        footer = ""
      } = content;

      Swal.fire({
        icon: "success",
        title: title,
        text: msg,
        footer: footer
      });
    }

    let error = function(content) {
      const {
        msg = "",
        title = "",
        footer = ""
      } = content;

      Swal.fire({
        icon: "error",
        title: title,
        text: msg,
        footer: footer
      });
    }

    let customHtmlAlert = async function(content) {
      const {
        msg = "",
        title = ""
      } = content;

      const { value: formValues } = await Swal.fire({
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showCancelButton: true,
        willOpen: () => {
          const reservationDatesModal = document.getElementById("reservation-dates-modal");
          console.log(reservationDatesModal);
          const dateRangePicker = new DateRangePicker(reservationDatesModal, {
            minDate: new Date(),
          });
        },
        preConfirm: () => {
          return [
            document.getElementById("check_in_date").value,
            document.getElementById("check_out_date").value
          ];
        },
        didOpen: () => {
            document.getElementById("check_in_date").removeAttribute("disabled");
            document.getElementById("check_out_date").removeAttribute("disabled");
        }
      });
      if (formValues) {
        Swal.fire(JSON.stringify(formValues));
      }

    }

    return {
      toast: toast,
      success: success,
      error: error,
      customHtmlAlert: customHtmlAlert
    }
}