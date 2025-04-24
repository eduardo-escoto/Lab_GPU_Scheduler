mkdir ./gpu_request
cd gpu_request
sudo mv /home/eduardo/db_importer_linux ./

chmod +x ./db_importer_linux
echo 'export DATABASE_DSN='\''gpu_sched_manager:JMcAuley4146!@tcp(deepfreeze.ucsd.edu:5129)/gpu_scheduler?parseTime=true'\''' >> ~/.bashrc
source ~/.bashrc
./db_importer_linux -verbose -table=gpus -mode=insert