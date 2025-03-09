#error_count=0
#for i in {1..100}; do
#    go run main.go 2>/dev/null || ((error_count++))
#done
#
#echo "Total errors: $error_count"

for i in {1..100}; do
    go run main.go 
done