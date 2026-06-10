ALTER TABLE treino.exe_exercicio 
ADD COLUMN usu_tx_id UUID NULL,
ADD CONSTRAINT fk_exercicio_usuario FOREIGN KEY (usu_tx_id) REFERENCES auth.usu_usuario(usu_tx_id);

CREATE INDEX idx_exercicio_usu_id ON treino.exe_exercicio(usu_tx_id);