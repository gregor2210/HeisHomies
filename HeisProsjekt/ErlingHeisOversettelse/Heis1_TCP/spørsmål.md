	
    ------ Fjerne Kommentarer er disse første 4 spørsmålene ---------


    InitIOHandling.go:
    1. Fjerne disse? Linje 10 og 16
    //drv_stop := make(chan bool)
    // go PollStopButton(drv_stop) // Ikke implementert enda.


    ElevatorIO.go:  `LØST`
    2. Fjerne denne? Linje 10
    // hvis denne endres, husk å endre i timer.go også


    configFSM.go: `LØST`
    3. Fjerne denne? Linje 5
    // IF THIS CHANGES, REMEMBER TO UPDATE IT IN ELEVATOR:IO.GO AS WELL

    main.go: `LØST`
    4. Skjønner ikke denne, ikke nødvendigvis fjerne. Linje 51
    // Make sure the elevator start running, set a call on current floor

    ------- Slette filer vi ikke trenger er spørsmål 4 til 6 ---------

    4. Slette main backup before chainge filen (Til innleveringen)

    5. Slette bugs.md filen (TIl innleveringen fjærner vi denne!)

    6. Slette cab_requests.txt inni fsm `BEHOLD DENNE!!!!!!!!!!!!!!!!!!!!!` (den kommer uansett tilbake, programmet vil bare generere en ny)

    7. Endre connectToElevatorserver() til connectToElevatorServer()?

    8. Endre elevio.InitIOHandling() til elevio.IOHandlingSetup() ? 

    9. Main L110: Endre errorBool til motorErrorDetectedBool eller bare motorErrorDetected?

    


---------