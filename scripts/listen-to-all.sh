#!/bin/bash
mosquitto_sub -h sphero.local -v \
  -t 'spheroRallye/roundEnd' \
  -t 'spheroRallye/player1/plannedActions' \
  -t 'spheroRallye/player1/possibleActionTypes' \
  -t 'spheroRallye/player1/errors' \
  -t 'spheroRallye/player1/adHocAction' 
  -t 'spheroRallye/player2/plannedActions' \
  -t 'spheroRallye/player2/possibleActionTypes' \
  -t 'spheroRallye/player2/errors' \
  -t 'spheroRallye/player2/adHocAction' 
  -t 'spheroRallye/player3/plannedActions' \
  -t 'spheroRallye/player3/possibleActionTypes' \
  -t 'spheroRallye/player3/errors' \
  -t 'spheroRallye/player3/adHocAction' 
  -t 'spheroRallye/player4/plannedActions' \
  -t 'spheroRallye/player4/possibleActionTypes' \
  -t 'spheroRallye/player4/errors' \
  -t 'spheroRallye/player4/adHocAction' 
  -t 'spheroRallye/player5/plannedActions' \
  -t 'spheroRallye/player5/possibleActionTypes' \
  -t 'spheroRallye/player5/errors' \
  -t 'spheroRallye/player5/adHocAction' 