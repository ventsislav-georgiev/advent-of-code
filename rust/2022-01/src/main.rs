use aoc::exec;

fn main() {
    exec(&[task1, task2]);
}

fn task1(input: String) {
    let mut max = 0;
    let mut cur = 0;

    for line in input.lines() {
        if line == "" {
            cur = 0;
            continue;
        }

        let val = line.parse::<i32>().unwrap();
        cur += val;

        if cur > max {
            max = cur
        }
    }

    println!("{}", max);
}

fn task2(input: String) {
    let mut vec: Vec<i32> = Vec::new();
    let mut max = 0;
    let mut cur = 0;

    for line in input.lines() {
        if line == "" {
            vec.push(cur);
            cur = 0;
            continue;
        }

        let val = line.parse::<i32>().unwrap();
        cur += val;

        if cur > max {
            max = cur
        }
    }

    vec.push(cur);
    vec.sort();
    let len = vec.len();

    println!("{}", vec[len - 1] + vec[len - 2] + vec[len - 3]);
}
