#include <thread>
#include <atomic>
#include <iostream>
#include <functional>

// g++ -std=c++11 main.cpp -o main && ./main && rm main

#define UNLOCKED false
#define LOCKED   true

class MutexImplicit
{
public:
    void Lock()
    {
        // std::memory_order order = std::memory_order_seq_cst
        while (m_state.exchange(LOCKED)) { 
            // итерация за итерацией
        }
    }

    void Unlock()
    {
        // std::memory_order order = std::memory_order_seq_cst
        m_state.store(UNLOCKED); 
    }

private:
    std::atomic<bool> m_state{ false };
};

class MutexExplicit
{
public:
    void Lock()
    {
        while (m_state.exchange(LOCKED, std::memory_order_acquire)) {
            // итерация за итерацией
        }
    }

    void Unlock()
    {
        m_state.store(UNLOCKED, std::memory_order_release);
    }

private:
    std::atomic<bool> m_state{ false };
};

template<typename MutexType>
void Benchmark(int threads_count)
{
    int counter = 0;
    MutexType mutex;

    std::function<void()> thread_fn = [&counter, &mutex]
    {
        for (int i = 0; i < 10; ++i) {
            mutex.Lock();
            ++counter;
            mutex.Unlock();
        }
    };

    std::vector<std::thread> threads(threads_count);
    for (int idx = 0; idx < threads.size(); ++idx)
        threads[idx] = std::thread(thread_fn);

    for (int idx = 0; idx < threads.size(); ++idx)
        threads[idx].join();
}

int main()
{   
    {
        std::chrono::steady_clock::time_point start = std::chrono::steady_clock::now();
        Benchmark<MutexImplicit>(2000);
        std::chrono::steady_clock::time_point finish = std::chrono::steady_clock::now();
        std::cout << "MutexImplicit = " << std::chrono::duration_cast<std::chrono::milliseconds>(finish - start).count() << "[ms]" << std::endl;
    }
    {
        std::chrono::steady_clock::time_point start = std::chrono::steady_clock::now();
        Benchmark<MutexExplicit>(2000);
        std::chrono::steady_clock::time_point finish = std::chrono::steady_clock::now();
        std::cout << "MutexExplicit = " << std::chrono::duration_cast<std::chrono::milliseconds>(finish - start).count() << "[ms]" << std::endl;
    }
    return EXIT_SUCCESS;
}