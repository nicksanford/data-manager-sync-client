# Build
```
make
```

# Usage

```bash
echo '{"address":"YOUR_VIAM_ROBOT_ADDRESS","id":"YOUR_ROBOT_PART_ID","secret":"YOUR_ROBOT_PART_API_SECRET"}' | jq > ./config.json

# to call [viam.service.datamanager.v1.DataManagerService.Sync](https://github.com/viamrobotics/api/blob/main/proto/viam/service/datamanager/v1/data_manager.proto#L15) every 5 seconds
./bin/PLATFORM/data-manager-sync-client ./config.json 5s

# to call [viam.service.datamanager.v1.DataManagerService.Sync](https://github.com/viamrobotics/api/blob/main/proto/viam/service/datamanager/v1/data_manager.proto#L15) once
./bin/PLATFORM/data-manager-sync-client ./config.json 
```
