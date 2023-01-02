use std::env;

pub fn exec(tasks: &[fn(String)]) {
    let args = env::args().collect();
    let year = get_int_arg(&args, "--year");
    let day = get_int_arg(&args, "--day");
    let task_number = get_int_arg(&args, "--task");

    let task = tasks[task_number as usize - 1];
    let input = get_input(year, day);
    task(input);
}

fn get_input(year: i64, day: i64) -> String {
    let session_cookie = "session=".to_owned() + env::var("SESSION_KEY").unwrap().as_str();
    let client = reqwest::blocking::Client::new();
    let res = client
        .get(format!(
            "https://adventofcode.com/{}/day/{}/input",
            year, day
        ))
        .header("COOKIE", session_cookie)
        .send()
        .unwrap();

    return res.text().unwrap();
}

fn get_int_arg(args: &Vec<String>, arg: &str) -> i64 {
    let mut result = 1;
    for a in args {
        if a.contains(arg) {
            result = a.split("=").collect::<Vec<&str>>()[1]
                .parse::<i64>()
                .unwrap();
        }
    }

    return result;
}
