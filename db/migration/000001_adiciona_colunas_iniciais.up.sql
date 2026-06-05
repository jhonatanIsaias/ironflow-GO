CREATE SCHEMA IF NOT EXISTS treino;

CREATE TABLE IF NOT EXISTS treino.exe_exercicio (
    exe_nr_id SERIAL PRIMARY KEY,
    exe_tx_nome VARCHAR(50) NOT NULL,
    exe_tx_grupo_muscular VARCHAR(50),
    exe_tx_grupo_muscular_sinegista VARCHAR(50),
    exe_tx_tipo_equipamento VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS treino.tre_treino (
    tre_nr_id SERIAL PRIMARY KEY,
    tre_tx_nome VARCHAR(60) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS treino.fit_ficha_treino (
    fit_nr_id SERIAL PRIMARY KEY,
    tre_nr_id INT NOT NULL,
    exe_nr_id INT NOT NULL,
    fit_nr_ordem INT DEFAULT 1,
    fit_nr_meta_series INT4,
    fit_nr_meta_repeticoes INT4,
    fit_nr_meta_peso DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_ficha_treino FOREIGN KEY (tre_nr_id) REFERENCES treino.tre_treino(tre_nr_id),
    CONSTRAINT fk_ficha_exercicio FOREIGN KEY (exe_nr_id) REFERENCES treino.exe_exercicio(exe_nr_id)
);

CREATE TABLE IF NOT EXISTS treino.set_sessao_treino (
    set_nr_id SERIAL PRIMARY KEY,
    tre_nr_id INT NOT NULL,
    set_dt_data DATE DEFAULT CURRENT_DATE,
    set_tm_hora_inicio TIME DEFAULT CURRENT_TIME,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_sessao_treino FOREIGN KEY (tre_nr_id) REFERENCES treino.tre_treino(tre_nr_id)
);

CREATE TABLE IF NOT EXISTS treino.sex_serie_executada (
    sex_nr_id SERIAL PRIMARY KEY,
    set_nr_id INT NOT NULL,
    fit_nr_id INT NOT NULL,
    sex_nr_serie_numero INT,
    sex_nr_repeticoes_realizadas INT4,
    sex_nr_peso_utilizado DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_serie_sessao FOREIGN KEY (set_nr_id) REFERENCES treino.set_sessao_treino(set_nr_id),
    CONSTRAINT fk_serie_ficha FOREIGN KEY (fit_nr_id) REFERENCES treino.fit_ficha_treino(fit_nr_id)
);