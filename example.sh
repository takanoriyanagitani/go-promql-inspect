#!/bin/bash

bin="./cmd/qlinspect/qlinspect"

export ENV_PROMQL_MAX_SIZE=65536

# 1. Simple Vector Selector
ex1(){
	echo "Testing: Simple Vector Selector"
	echo 'load_avg1' | "${bin}"
	echo "-------------------"
}

# 2. Range Vector with Function and Label Matcher
ex2(){
	echo "Testing: Rate function with label matcher and range"
	echo 'rate(cpu_ms{job="job1"}[5m])' | "${bin}"
	echo "-------------------"
}

# 3. Complex Label Matchers (Regex and Negative)
ex3(){
	echo "Testing: Regex (=~) and Negative (!=) matchers"
	echo 'up{instance=~"node-.*", job!="prometheus"}' | "${bin}"
	echo "-------------------"
}

# 4. Nested Functions (Aggregation of a Function)
ex4(){
	echo "Testing: Nested functions (sum of rates)"
	echo 'sum(rate(http_requests_total{status="200"}[5m]))' | "${bin}"
	echo "-------------------"
}

# 5. Arithmetic Operations (Tests the 'default' node logger)
ex5(){
	echo "Testing: Arithmetic expression (+ and /)"
	echo '(node_memory_MemTotal_bytes - node_memory_MemFree_bytes) / node_memory_MemTotal_bytes' | "${bin}"
	echo "-------------------"
}

# 6. Range Selector with Offset
ex6(){
	echo "Testing: Range selector with offset"
	echo 'node_cpu_seconds_total{mode="idle"}[1h] offset 2d' | "${bin}"
	echo "-------------------"
}

# 7. Aggregation by Group (sum by)
ex7(){
	echo "Testing: Aggregation with grouping labels"
	echo 'sum by (instance, job) (up)' | "${bin}"
	echo "-------------------"
}

# 8. Complex Logical Expression (and/or/unless)
ex8(){
	echo "Testing: Logical binary operator"
	echo 'node_cpu_seconds_total > 0.1 and node_load1 > 5.3' | "${bin}"
	echo "-------------------"
}

#ex1
#ex2
#ex3
#ex4
#ex5
#ex6
#ex7
ex8
