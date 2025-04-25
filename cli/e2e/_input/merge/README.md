# Input Setup
```bash
cd ./cli/

# current.sqlite
rm -f ./e2e/_input/merge/current.sqlite
./hyaline --debug extract current --config ./e2e/_input/merge/config-current.yml --system git --output ./e2e/_input/merge/current.sqlite

# current-copy.sqlite
rm -f ./e2e/_input/merge/current-copy.sqlite
cp ./e2e/_input/merge/current.sqlite ./e2e/_input/merge/current-copy.sqlite
# Open database and make the following changes...
open ./e2e/_input/merge/current-copy.sqlite
# DELETE FROM CODE WHERE ID="app-http";
# DELETE FROM FILE WHERE CODE_ID="app-http";
# DELETE FROM DOCUMENTATION WHERE ID="app-http";
# DELETE FROM DOCUMENT WHERE DOCUMENTATION_ID="app-http";
# DELETE FROM SECTION WHERE DOCUMENTATION_ID="app-http";
# UPDATE CODE SET PATH="../../hyaline-example-copy/" WHERE ID="app-path";
# UPDATE FILE SET RAW_DATA=CONCAT(RAW_DATA, " // Now with copy") WHERE CODE_ID="app-path" AND ID="index.js";
# UPDATE DOCUMENTATION SET PATH="../../hyaline-example-copy/" WHERE ID="app-path";
# UPDATE DOCUMENT SET RAW_DATA=CONCAT(RAW_DATA, " // Now with copy") WHERE DOCUMENTATION_ID="app-path" AND ID="README.md";
# UPDATE DOCUMENT SET EXTRACTED_DATA=CONCAT(EXTRACTED_DATA, " // Now with copy") WHERE DOCUMENTATION_ID="app-path" AND ID="README.md";
# UPDATE SECTION SET EXTRACTED_DATA=CONCAT(EXTRACTED_DATA, " - Now with copy") WHERE DOCUMENTATION_ID="app-path" AND ID="README.md";
# UPDATE SECTION SET EXTRACTED_DATA=CONCAT(EXTRACTED_DATA, " - Now with copy") WHERE DOCUMENTATION_ID="app-path" AND ID="README.md#Example#Subsection 2";

# change.sqlite
rm -f ./e2e/_input/merge/change.sqlite
./hyaline --debug extract change --config ./e2e/_input/merge/config-change.yml --system git --base main --head origin/feat-1 --pull-request appgardenstudios/hyaline-example/1 --issue appgardenstudios/hyaline-example/2 --issue appgardenstudios/hyaline-example/3  --output ./e2e/_input/merge/change.sqlite
```