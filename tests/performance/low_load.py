from locust import LoadTestShape

class LowLoad(LoadTestShape):
    time_limit = 60
    spawn_rate = 2
    max_users = 10

    def tick(self) -> (tuple[float, int] | None):
        user_count = self.spawn_rate * self.get_run_time()
        if user_count > self.max_users:
            user_count = self.max_users

        return (user_count, self.spawn_rate)