# Lab-2-SD
Laboratorio numero 2 de sistemas distribuidos Golang-grpc  
Gianni Carlini 201773105-2

Máquina 1 (10.10.28.67) - DataNode 1  
Máquina 2 (10.10.28.68) - DataNode 2  
Máquina 3 (10.10.28.69) - DataNode 3  
Máquina 4 (10.10.28.7)  - NameNode  
 
El sistema se incia en las maquinas designadas arriba con el make data1/data2/data3 para los datanodes, make cliente para el cliente en cualquira de las maquinas y make namenode en la maquina 4. Por parte del cliente primero debes decidir cual sera el cordinador y elegir el sistema a simular, habiendo elegido esto te permitira subir o descargar archivos, primero se debe subir un archivo cual sea (Dracula, Frankenstein, MobyDick, Mujercitas, PeterPan) para llenar el sistema, luego puedes descargar.
Al iniciar los otros sistemas (datas y name) se pedira elegir el sistema a simular pudiendo ingresar 1 o 2 (centralizado/distribuido) en el sistema centralizado, debes informarle a los Datanodes si es o no el cordinador con un 1 si lo es o 2 si no lo es.   
Para probar los distintos sistemas es necesario elegir el mismo modo de ejecucion entodos los nodos y cliente, para correr el otro sistema cancelar la ejecucion y volver a ejercurtalo.  

En caso de necesitar recompilar el proto ejecutar make grpc.
